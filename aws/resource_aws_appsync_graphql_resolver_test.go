package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSAppsyncResolver_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsAppsyncGraphqlResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncGraphqlResolverConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppsyncGraphqlResolverExists("aws_appsync_graphql_resolver.test"),
					resource.TestCheckResourceAttrSet("aws_appsync_graphql_resolver.", "arn"),
				),
			},
		},
	})
}

func testAccCheckAwsAppsyncGraphqlResolverExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*AWSClient).appsyncconn
		input := &appsync.GetResolverInput{
			ApiId: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetResolver(input)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAwsAppsyncGraphqlResolverDestroy(s *terraform.State) error {
	return nil
}

func testAccAppsyncGraphqlResolverConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name = "tf_appsync_%s"
}

resource "aws_elasticsearch_domain" "test" {
  domain_name = "tf-es-%s"
  ebs_options {
    ebs_enabled = true
    volume_size = 10
  }
}

resource "aws_iam_role" "test" {
  name = "tf-role-%s"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "appsync.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  name = "tf-rolepolicy-%s"
  role = "${aws_iam_role.test.id}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "es:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_elasticsearch_domain.test.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = "${aws_appsync_graphql_api.test.id}"
  name = "tf_appsync_%s"
  type = "AMAZON_ELASTICSEARCH"
  elasticsearch_config {
    region = "${data.aws_region.current.name}"
    endpoint = "https://${aws_elasticsearch_domain.test.endpoint}"
  }
  service_role_arn = "${aws_iam_role.test.arn}"
}

resource "aws_appsync_graphql_resolver" "test" {
	api_id = "${aws_appsync_graphql_api.test.id}"
	data_source_name = "${aws_appsync_datasource.test.name}"
	type_name = "Query"
	field_name = "getPost"
	request_mapping_template = <<EOF
{
    "version":"2017-02-28",
    "operation":"GET",
    "path":"/id/post/_search",
    "params":{
        "headers":{},
        "queryString":{},
        "body":{
            "from":0,
            "size":50
        }
    }
}
EOF
	response_mapping_template = <<EOF
[
    #foreach($entry in $context.result.hits.hits)
    #if( $velocityCount > 1 ) , #end
    $utils.toJson($entry.get("_source"))
    #end
]
EOF
}
`, rName, rName, rName, rName, rName)
}

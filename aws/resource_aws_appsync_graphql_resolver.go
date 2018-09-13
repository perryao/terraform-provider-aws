package aws

//https://docs.aws.amazon.com/appsync/latest/devguide/resolver-mapping-template-reference-overview.html

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsAppsyncGraphqlResolver() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsAppsyncGraphqlResolverCreate,
		Read:   resourceAwsAppsyncGraphqlResolverRead,
		Update: resourceAwsAppsyncGraphqlResolverUpdate,
		Delete: resourceAwsAppsyncGraphqlResolverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_mapping_template": {
				Type:     schema.TypeString,
				Required: true,
			},
			// TODO: the AWS api seems to require this
			"response_mapping_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsAppsyncGraphqlResolverCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appsyncconn

	input := &appsync.CreateResolverInput{
		ApiId:                  aws.String(d.Get("api_id").(string)),
		DataSourceName:         aws.String(d.Get("data_source_name").(string)),
		FieldName:              aws.String(d.Get("field_name").(string)),
		TypeName:               aws.String(d.Get("type_name").(string)),
		RequestMappingTemplate: aws.String(d.Get("request_mapping_template").(string)),
	}

	if v, ok := d.GetOk("response_mapping_template"); ok {
		input.ResponseMappingTemplate = aws.String(v.(string))
	}

	resp, err := conn.CreateResolver(input)
	if err != nil {
		return err
	}

	d.Set("arn", resp.Resolver.ResolverArn)

	return nil
}

func resourceAwsAppsyncGraphqlResolverRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsAppsyncGraphqlResolverUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsAppsyncGraphqlResolverDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

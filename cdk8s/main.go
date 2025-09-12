package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus32/v2"
)

func NewBot(scope constructs.Construct, id string) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), nil)

	label := map[string]*string{"app": jsii.String("aegisbot")}
	tokenSecret := cdk8splus32.Secret_FromSecretName(chart, jsii.String("Secret"), jsii.String("aegisbot"))

	// StatefulSet later?
	cdk8splus32.NewDeployment(chart, jsii.String("aegisbot"), &cdk8splus32.DeploymentProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Labels: &label,
		},
		Replicas: jsii.Number(1),
		Containers: &[]*cdk8splus32.ContainerProps{
			{
				Name:  jsii.String("bot"),
				Image: jsii.String("kkrypt0nn/aegisbot:latest"),
				EnvVariables: &map[string]cdk8splus32.EnvValue{
					"BOT_TOKEN": cdk8splus32.EnvValue_FromSecretValue(&cdk8splus32.SecretValue{
						Key:    jsii.String("BOT_TOKEN"),
						Secret: tokenSecret,
					}, nil),
				},
			},
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewBot(app, "aegisbot")
	app.Synth()
}

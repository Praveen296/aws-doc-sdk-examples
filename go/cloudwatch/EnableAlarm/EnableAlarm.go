/*
   Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/
// snippet-start:[cloudwatch.go.enable_alarm]
package main

// snippet-start:[cloudwatch.go.enable_alarm.imports]
import (
    "flag"
    "fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
)
// snippet-end:[cloudwatch.go.enable_alarm.imports]

// EnableAlarm creates an alarm that reboots an EC2 instance whenever the CPU utilization is greater than 70%
// Inputs:
//     instanceName is the name of the EC2 instance
//     instanceID is the ID of the EC2 instance
//     alarmName is the name of the alarm
// Output:
//     If successful, nil
//     Otherwise, the error from a call to PutMetricAlarm or EnableAlarmActions
func EnableAlarm(instanceName string, instanceID string, alarmName string) error {

    // Initialize a session that the SDK uses to load
    // credentials from the shared credentials file ~/.aws/credentials
    // and configuration from the shared configuration file ~/.aws/config.
    // snippet-start:[cloudwatch.go.enable_alarm.session]
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create new CloudWatch client
    svc := cloudwatch.New(sess)
    // snippet-end:[cloudwatch.go.enable_alarm.session]
    
    // snippet-start:[cloudwatch.go.enable_alarm.put]
    // Get region for alarm action
    region := svc.Config.Region


    _, err := svc.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
        AlarmName:          aws.String(alarmName),
        ComparisonOperator: aws.String(cloudwatch.ComparisonOperatorGreaterThanThreshold),
        EvaluationPeriods:  aws.Int64(1),
        MetricName:         aws.String("CPUUtilization"),
        Namespace:          aws.String("AWS/EC2"),
        Period:             aws.Int64(60),
        Statistic:          aws.String(cloudwatch.StatisticAverage),
        Threshold:          aws.Float64(70.0),
        ActionsEnabled:     aws.Bool(true),
        AlarmDescription:   aws.String("Alarm when server CPU exceeds 70%"),
        Unit:               aws.String(cloudwatch.StandardUnitSeconds),
        AlarmActions: []*string{
            aws.String(fmt.Sprintf("arn:aws:swf:"+*region+":%s:action/actions/AWS_EC2.InstanceId.Reboot/1.0", instanceName)),
        },
        Dimensions: []*cloudwatch.Dimension{
            {
                Name:  aws.String("InstanceId"),
                Value: aws.String(instanceID),
            },
        },
    })
    // snippet-end:[cloudwatch.go.enable_alarm.put]
    if err != nil {
        return err
    }

    // Enable the alarm for the instance
    // snippet-start:[cloudwatch.go.enable_alarm.enable]
    _, err = svc.EnableAlarmActions(&cloudwatch.EnableAlarmActionsInput{
        AlarmNames: []*string{
            aws.String(instanceID),
        },
    })
    // snippet-end:[cloudwatch.go.enable_alarm.enable]
    if err != nil {
        return err
    }

    return nil
}

func main() {
    instanceNamePtr := flag.String("n", "", "The instance name")
    instanceIDPtr := flag.String("i", "", "The instance ID")
    alarmNamePtr := flag.String("a", "", "The alarm name")
    flag.Parse()
    instanceName := *instanceNamePtr
    instanceID := *instanceIDPtr
    alarmName := *alarmNamePtr

    if instanceName == "" || instanceID == "" || alarmName == "" {
        fmt.Println("You must supply an instance name, instance ID, and alarm name")
        return
    }

    err := EnableAlarm(instanceName, instanceID, alarmName)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Enabled alarm " + alarmName + " for EC2 instance " + instanceName)
}
// snippet-end:[cloudwatch.go.enable_alarm]
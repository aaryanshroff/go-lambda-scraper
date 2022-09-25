terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

data "archive_file" "scraper" {
  type        = "zip"
  source_file = "main"
  output_path = "scraper.zip"
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : "sts:AssumeRole",
          "Principal" : {
            "Service" : "lambda.amazonaws.com"
          },
          "Effect" : "Allow",
          "Sid" : ""
        }
      ]
  })

  managed_policy_arns = [aws_iam_policy.dynamodb_for_lambda.arn]
}

resource "aws_iam_policy" "dynamodb_for_lambda" {
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "ReadWriteTable",
        "Effect" : "Allow",
        "Action" : [
          "dynamodb:BatchGetItem",
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:BatchWriteItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem"
        ],
        "Resource" : "arn:aws:dynamodb:*:*:table/ScraperHistory"
      },
      {
        "Sid" : "GetStreamRecords",
        "Effect" : "Allow",
        "Action" : "dynamodb:GetRecords",
        "Resource" : "arn:aws:dynamodb:*:*:table/ScraperHistory/stream/* "
      },
      {
        "Sid" : "WriteLogStreamsAndGroups",
        "Effect" : "Allow",
        "Action" : [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        "Resource" : "*"
      },
      {
        "Sid" : "CreateLogGroup",
        "Effect" : "Allow",
        "Action" : "logs:CreateLogGroup",
        "Resource" : "*"
      }
    ]
  })
}

resource "aws_lambda_function" "scraper" {
  filename      = "scraper.zip"
  function_name = "Scraper"
  handler       = "main"
  role          = aws_iam_role.iam_for_lambda.arn

  source_code_hash = data.archive_file.scraper.output_base64sha256

  runtime = "go1.x"
}

resource "aws_dynamodb_table" "scraper-history" {
  name           = "ScraperHistory"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "URL"

  attribute {
    name = "URL"
    type = "S"
  }

}

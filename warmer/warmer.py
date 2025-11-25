import boto3


def warmer(event, context):
    client = boto3.client("lambda")

    paginator = client.get_paginator('list_functions')
    for page in paginator.paginate(PaginationConfig={
        'PageSize': 50
    }):
        for function in page['Functions']:
            if 'parking-management-backend' in function['FunctionName'] and function['FunctionName'] != context.function_name:
                print(f"warming {function['FunctionName']}")
                client.invoke(
                    FunctionName=function['FunctionName'],
                    InvocationType='Event'
                )

# warmer(None, None)
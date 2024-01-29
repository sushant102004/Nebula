/*
	@author: Sushant
	@last-modified: 23 January 2024
	@GitHub: https://github.com/sushant102004
*/

import { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda'

export const handler = async (event : APIGatewayProxyEvent) : Promise<APIGatewayProxyResult> => {
    try {
        return {
            body: JSON.stringify({
                message : 'Function Called'
            }),
            statusCode: 200
        }
    } catch (err) {
        return {
            body: JSON.stringify({
                error : err,
            }),
            statusCode: 500
        }
    }
}
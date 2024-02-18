/*
    @author: Sushant
    @created: 11 February 2024
    @GitHub: https://github.com/sushant102004

    PoF: - Provide functionality for image compressing and resizing.

    Request: - Accept base64 encoded image data and quality. Lesser the quality is more the compression will be applied.
*/

import { APIGatewayProxyResult, APIGatewayProxyEvent } from 'aws-lambda'
import sharp from 'sharp'

interface RequestBody {
    image: string
    quality: string
}

export const handler = async (event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> => {
    try {
        if (!event.body) {
            return {
                statusCode: 400,
                body: JSON.stringify({   
                    message: 'Provide valid input data.'
                })
            }
        }

        const eventBody : RequestBody = JSON.parse(event.body)

        const imageContent = eventBody.image;

        const inputImageBuffer = base64ToArrayBuffer(imageContent);


        const imageBuffer = await sharp(inputImageBuffer.buffer).webp({
            quality: parseInt(eventBody.quality)
        }).toBuffer();


        return {
            statusCode: 200,
            headers: {
                'Content-Type': 'image/webp'
            },
            body: JSON.stringify({
                resizedImage: imageBuffer.toString('base64'),
                message : 'OK'
            })
        };
    } catch (err: any) {
        return {
            statusCode: 500,
            body: JSON.stringify({
                message: 'Some error happened',
                error: err.stack
            }),
        };
    }
};

function base64ToArrayBuffer(base64: string): Uint8Array {
    let binaryString = atob(base64)
    let bytes = new Uint8Array(base64.length)

    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes
}
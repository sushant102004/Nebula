/*
    @author: Sushant
    @created: 11 February 2024
    @GitHub: https://github.com/sushant102004

    PoF: - Provide functionality for image compressing and resizing.
*/

import { APIGatewayProxyResult } from 'aws-lambda'
import sharp from 'sharp'

type ResizeImageEvent = {
    image: string
    quality: number
}


export const handler = async (event : ResizeImageEvent) : Promise<APIGatewayProxyResult> => {
    try {
        const imageContent = event.image
        
        const inputImageBuffer = base64ToArrayBuffer(imageContent)


        const imageBuffer = await sharp(inputImageBuffer.buffer).jpeg({
            quality: event.quality
        }).toBuffer()
    

        return {
            statusCode: 200,
            body: JSON.stringify({
                resizedImage: imageBuffer.toString('base64')
            })
        }
    } catch (err : any) {
        return {
            statusCode: 500,
            body: JSON.stringify({
                message: 'some error happened',
                error : err.stack
            }),
        };
    }
}

function base64ToArrayBuffer(base64 : string) : Uint8Array {
    let binaryString = atob(base64)
    let bytes = new Uint8Array(base64.length)

    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes
}
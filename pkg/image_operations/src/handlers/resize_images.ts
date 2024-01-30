/*
    @author: Sushant
    @created: 30 January 2024
    @last-modified: 30 January 2024
    @GitHub: https://github.com/sushant102004

    PoF: - Provide functionality for image resizing.

    Working Steps: -
    1. Get base64 encoding of image in request. -> MVP
    2. Apply resizing functionality. -> MVP
    3. Return new image base64 encoding. -> MVP
*/

import { Handler, APIGatewayProxyResult } from '../../dist/node_modules/@types/aws-lambda'
import sharp from '../../dist/node_modules/sharp/lib'

type ResizeImageEvent = {
    // base64 encoded image
    image: string
    width: number
    height: number
    fit: sharp.FitEnum

    // Background color when using fit of 'contain'
    // background: string
}


export const handler: Handler = async (event: ResizeImageEvent): Promise<APIGatewayProxyResult> => {
    try {
        const imageContent = event.image
        
        const inputImageBuffer = base64ToArrayBuffer(imageContent)

        const imageBuffer = await sharp(inputImageBuffer.buffer)
            .resize({
                width: event.width,
                height: event.height
            }).toBuffer()
    

        return {
            statusCode: 200,
            body: JSON.stringify({
                resizedImage: imageBuffer.toString('base64')
            })
        }
    } catch (err) {
        return {
            statusCode: 500,
            body: JSON.stringify({
                message: 'some error happened',
                error : err
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
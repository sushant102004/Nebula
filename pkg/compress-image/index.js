"use strict";
/*
    @author: Sushant
    @created: 11 February 2024
    @GitHub: https://github.com/sushant102004

    PoF: - Provide functionality for image compressing and resizing.
*/
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.handler = void 0;
const sharp_1 = __importDefault(require("sharp"));
const handler = (event) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        if (!event.body) {
            return {
                statusCode: 400,
                body: JSON.stringify({
                    message: 'Provide valid input data.'
                })
            };
        }
        const eventBody = JSON.parse(event.body);
        const imageContent = eventBody.image;
        const inputImageBuffer = base64ToArrayBuffer(imageContent);
        const imageBuffer = yield (0, sharp_1.default)(inputImageBuffer.buffer).webp({
            quality: parseInt(eventBody.quality)
        }).toBuffer();
        return {
            statusCode: 200,
            headers: {
                'Content-Type': 'image/webp'
            },
            body: JSON.stringify({
                resizedImage: imageBuffer.toString('base64'),
                message: 'OK'
            })
        };
    }
    catch (err) {
        return {
            statusCode: 500,
            body: JSON.stringify({
                message: 'Some error happened',
                error: err.stack
            }),
        };
    }
});
exports.handler = handler;
function base64ToArrayBuffer(base64) {
    let binaryString = atob(base64);
    let bytes = new Uint8Array(base64.length);
    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
}

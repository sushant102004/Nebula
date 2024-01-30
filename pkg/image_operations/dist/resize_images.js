"use strict";
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
const lib_1 = __importDefault(require("./node_modules/sharp/lib"));
const handler = (event) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const imageContent = event.image;
        const inputImageBuffer = base64ToArrayBuffer(imageContent);
        const imageBuffer = yield (0, lib_1.default)(inputImageBuffer.buffer)
            .resize({
            width: event.width,
            height: event.height
        }).toBuffer();
        return {
            statusCode: 200,
            body: JSON.stringify({
                resizedImage: imageBuffer.toString('base64')
            })
        };
    }
    catch (err) {
        return {
            statusCode: 500,
            body: JSON.stringify({
                message: 'some error happened',
                error: err
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

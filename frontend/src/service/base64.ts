export default function base64UrlToString(base64Url: string): string {
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    while (base64.length % 4) {
        base64 += '=';
    }

    const binaryString = atob(base64);
    const bytes = Uint8Array.from(binaryString, char => char.charCodeAt(0));

    return new TextDecoder().decode(bytes);
}
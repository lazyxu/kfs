export function modeIsDir(mode) {
    return mode >= 2147483648;
}

export function getPerm(mode) {
    return mode >= 2147483648 ? mode - 2147483648 : 0;
}

import moment from 'moment';

export function modeIsDir(mode) {
    return mode >= 2147483648;
}

export function getPerm(mode) {
    return mode >= 2147483648 ? mode - 2147483648 : mode;
}

/**
 * 将 int64 权限值转换为 drwxr-xr-x 格式字符串
 * @param {bigint|number} mode - int64 权限值（支持 number 或 bigint）
 * @returns {string} 权限字符串（如 'drwxr-xr-x'）
 */
export function formatLinuxPermission(mode) {
    // 处理 int64 的 32 位限制：仅取低 16 位（高 4 位文件类型 + 低 12 位权限）
    const safeMode = BigInt(mode) & 0xFFFFn; // 限制在 16 位范围内[5,7](@ref)
    
    // 1. 提取文件类型（高 4 位）
    const fileTypeBits = Number(safeMode >> 12n) & 0xF;
    let typeChar = '-'; // 默认普通文件
    switch (fileTypeBits) {
        case 0x4: typeChar = 'd'; break; // 目录[4](@ref)
        case 0x8: typeChar = '-'; break; // 普通文件
        case 0xA: typeChar = 'l'; break; // 符号链接
        // 可扩展其他类型：管道(p)、套接字(s)等
    }

    // 2. 提取权限位（低 12 位，仅用低 9 位）
    const permBits = Number(safeMode & 0x1FFn); // 0x1FF = 二进制 111111111 (9位)[4](@ref)

    // 3. 权限分组转换函数
    const toRwx = (perm) => {
        const r = (perm & 4) ? 'r' : '-'; // 100 (读)
        const w = (perm & 2) ? 'w' : '-'; // 010 (写)
        const x = (perm & 1) ? 'x' : '-'; // 001 (执行)
        return `${r}${w}${x}`;
    };

    // 4. 拆分三组权限
    const owner = (permBits >> 6) & 0x7; // 高3位（7-9位）
    const group = (permBits >> 3) & 0x7; // 中3位（4-6位）
    const others = permBits & 0x7;       // 低3位（1-3位）[4](@ref)

    // 5. 组合结果
    return typeChar + toRwx(owner) + toRwx(group) + toRwx(others);
}

export function isDCIM(name) {
    name = name.toLowerCase()
    if (name.endsWith(".jpg") || name.endsWith(".jpeg") || name.endsWith(".png") || name.endsWith(".heic") ||
        name.endsWith(".mp4") || name.endsWith(".mov")
    ) {
        return true;
    }
    return false;
}

export function isImage(name) {
    name = name.toLowerCase()
    if (name.endsWith(".jpg") || name.endsWith(".jpeg") || name.endsWith(".png") || name.endsWith(".heic")) {
        return true;
    }
    return false;
}

export function isVideo(name) {
    name = name.toLowerCase()
    if (name.endsWith(".mp4") || name.endsWith(".mov")) {
        return true;
    }
    return false;
}

export function isViewable(name) {
    name = name.toLowerCase()
    if (name.endsWith(".txt") || name.endsWith(".md") || name.endsWith(".log") ||
        name.endsWith(".go") || name.endsWith(".js") || name.endsWith(".java") || name.endsWith(".c")
    ) {
        return true;
    }
    return false;
}

export function parseShotTime(metadata) {
    let { exif, videoMetadata } = metadata;
    if (exif) {
        if (exif.DateTime) {
            return moment.parseZone(exif.DateTime + " " + exif.OffsetTime, "YYYY:MM:DD HH:mm:ss ZZ");
        } else if (exif.DateTimeOriginal) {
            return moment.parseZone(exif.DateTimeOriginal + " " + exif.OffsetTimeOriginal, "YYYY:MM:DD HH:mm:ss ZZ");
        }
        return moment.parseZone(exif.DateTimeDigitized + " " + exif.OffsetTimeDigitized, "YYYY:MM:DD HH:mm:ss ZZ");
    } else if (videoMetadata) {
        if (videoMetadata.Created) {
            return moment(videoMetadata.Created / 1000 / 1000);
        }
        return moment(videoMetadata.Modified / 1000 / 1000);
    }
    return moment.invalid();
}

export function timeSortFn(a, b) {
    if (!a.shotTime.isValid()) {
        return -1;
    } else if (!b.shotTime.isValid()) {
        return 1;
    }
    return a.shotTime.isAfter(b.shotTime) ? 1 : -1;
}

export function parseShotEquipment(metadata) {
    let { exif } = metadata;
    if (exif) {
        if (exif.Model) {
            if (exif.Model.includes(exif.Make)) {
                return exif.Model;
            } else {
                return exif.Make + " " + exif.Model;
            }
        } else {
            return exif.HostComputer;
        }
    }
    return "";
}

export function toPrecent(n) {
    return (Math.floor(n * 100000) / 1000).toFixed(3) + "%";
}

export function getTransform(orientation) {
    switch (orientation) {
        case 2:
            return "rotateY(180deg)";
        case 3:
            return "rotate(180deg)";
        case 4:
            return "rotate(180deg)rotateY(180deg)";
        case 5:
            return "rotate(270deg)rotateY(180deg)";
        case 6:
            return "rotate(90deg)";
        case 7:
            return "rotate(90deg)rotateY(180deg)";
        case 8:
            return "rotate(270deg)";
    }
    return "";
}
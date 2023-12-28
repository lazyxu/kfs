import { parseShotEquipment } from "@kfs/common/api/utils";
import { getMetadata } from "@kfs/common/api/webServer/exif";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useState } from 'react';
import { Pressable, View } from 'react-native';
import { Appbar, Text } from "react-native-paper";

function formatTime(DateTime, SubsecTime, OffsetTime) {
    let t = "";
    if (DateTime) {
        t += DateTime;
    }
    if (SubsecTime) {
        t += "." + SubsecTime;
    }
    if (OffsetTime) {
        t += " " + OffsetTime;
    }
    return t;
}

function parseLocation(metadata) {
    let { exif } = metadata;
    if (exif) {
        let t = "";
        if (exif.GPSLatitudeRef) {
            t += exif.GPSLatitudeRef + exif.GPSLatitude + ", " + exif.GPSLongitudeRef + exif.GPSLongitude;
        }
        return t;
    }
}

export default function ({ navigation, route }) {
    const { hash } = route.params;
    const [metadata, setMetadata] = useState();
    const [sameFiles, setSameFiles] = useState([]);
    const sysConfig = getSysConfig();
    useEffect(() => {
        getMetadata(hash).then(setMetadata);
        // listDriverFileByHash(hash).then(sf => {
        //     setSameFiles(sf);
        // });
    }, []);
    console.log("Info", hash, metadata);
    return (
        <>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="信息" />
                <Pressable key={hash} onPress={() => navigation.pop()}  >
                    <Text>完成</Text>
                </Pressable>
            </Appbar.Header>
            {!metadata ?
                <View style={{
                    flex: 1,
                    alignItems: 'center',
                    justifyContent: 'center'
                }}>
                    <Text >TODO</Text>
                </View> : <>
                    <Text>哈希值：{hash}</Text>
                    {metadata.exif && <>
                        <Text>原始时间：{formatTime(metadata.exif.DateTimeOriginal, metadata.exif.SubsecTimeOriginal, metadata.exif.OffsetTimeOriginal)}</Text>
                        <Text>数字化时间：{formatTime(metadata.exif.DateTimeDigitized, metadata.exif.SubsecTimeDigitized, metadata.exif.OffsetTimeDigitized)}</Text>
                        <Text>修改时间：{formatTime(metadata.exif.DateTime, metadata.exif.SubsecTime, metadata.exif.OffsetTime)}</Text>
                        <Text>图片信息：{ }</Text>
                        <Text>相机信息：{parseShotEquipment(metadata)}</Text>
                        <Text>位置：{parseLocation(metadata)}</Text>
                    </>}
                </>}
        </>
    );
}

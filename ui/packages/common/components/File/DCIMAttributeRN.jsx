import { parseShotEquipment } from "@kfs/common/api/utils";
import { getMetadata } from "@kfs/common/api/webServer/exif";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useState } from 'react';
import { Pressable } from 'react-native';
import { Appbar, Surface, Text, Modal, Portal, IconButton, Divider } from "react-native-paper";
import { ScrollView, View } from 'react-native';
import humanize from 'humanize';

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
    let t = "";
    if (metadata.gps_latitude_ref) {
        t += metadata.gps_latitude_ref + metadata.gps_latitude + ", " + metadata.gps_longitude_ref + metadata.gps_longitude;
    }
    return t;
}

export default function ({ dcimFile, onClose, navigation, route }) {
    // console.log("DCIMAttribute.route.params", route.params);
    const { id, size, md5, sha256 } = dcimFile;
    const [metadata, setMetadata] = useState();
    const [sameFiles, setSameFiles] = useState([]);
    const [modalVisible, setModalVisible] = useState(false);
    const sysConfig = getSysConfig();
    useEffect(() => {
        getMetadata(id).then(setMetadata);
        // listDriverFileByHash(hash).then(sf => {
        //     setSameFiles(sf);
        // });
    }, []);
    console.log("DCIMAttribute", id, dcimFile, metadata);
    return (
        <Surface style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="信息" />
                {/* <Pressable onPress={() => navigation.pop()}  > */}
                <Pressable onPress={onClose}  >
                    <Text>完成</Text>
                </Pressable>
            </Appbar.Header>
            {!metadata ?
                <Surface style={{
                    flex: 1,
                    alignItems: 'center',
                    justifyContent: 'center'
                }}>
                    <Text >TODO</Text>
                </Surface>
                : <Surface style={{
                    flex: 1
                }}>
                    <Text>文件ID：{id}</Text>
                    <Text>文件大小：{humanize.filesize(size)}</Text>
                    <Text>md5：{md5}</Text>
                    <Text>sha256：{sha256}</Text>
                    <Text>原始时间：{formatTime(metadata.date_time_original, metadata.subsec_time_original, metadata.offset_time_original)}</Text>
                    <Text>数字化时间：{formatTime(metadata.date_time_digitized, metadata.subsec_time_digitized, metadata.offset_time_digitized)}</Text>
                    <Text>修改时间：{formatTime(metadata.date_time, metadata.subsec_time, metadata.offset_time)}</Text>
                    <Text>图片信息：{ }</Text>
                    <Text>相机信息：{parseShotEquipment(metadata)}</Text>
                    <Text>位置：{parseLocation(metadata)}</Text>
                </Surface>}
            </Surface>
    );
}

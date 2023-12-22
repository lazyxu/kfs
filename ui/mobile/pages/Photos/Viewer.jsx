import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useState } from "react";
import { Image, View } from 'react-native';
import { Appbar } from "react-native-paper";

export default function ({ navigation, route }) {
    const { hash } = route.params;
    console.log("Viewer", hash);
    const [downloadName, setDownloadName] = useState();
    const [metadata, setMetadata] = useState();
    const [openMetadata, setOpenMetadata] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const sysConfig = getSysConfig();
    useEffect(() => {
        // getMetadata(hash).then(setMetadata);
        // listDriverFileByHash(hash).then(sf => {
        //     setSameFiles(sf);
        //     setDownloadName(sf[0].name);
        // });
    }, []);
    return (
        <>
            <Appbar.Header>
                <Appbar.BackAction onPress={() => navigation.pop()} />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <View style={{
                padding: "0",
                width: "100%",
                height: "100%",
                display: "flex",
                justifyContent: "center",
                alignItems: "center"
            }}>
                <Image style={{
                    maxWidth: "100%", maxHeight: "100%",
                    width: "100%",
                    height: "100%",
                }} source={`${sysConfig.webServer}/api/v1/image?hash=${hash}`} />
            </View>
        </>
    );
}

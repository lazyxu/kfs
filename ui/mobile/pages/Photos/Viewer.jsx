import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useState } from "react";
import { Image } from 'react-native';
import { Appbar } from "react-native-paper";

export default function ({ hash, onClose }) {
    console.log(hash, onClose);
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
                <Appbar.BackAction onPress={onClose} />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <Image style={{ maxWidth: "100%", maxHeight: "100%" }} source={`${sysConfig.webServer}/api/v1/image?hash=${hash}`} />
        </>
    );
}

import { downloadURI } from "@kfs/common/api/web";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import * as Sharing from 'expo-sharing';
import React, { useEffect, useRef, useState } from 'react';
import { Image, Platform } from 'react-native';
import { Appbar } from "react-native-paper";
import ViewShot from 'react-native-view-shot';

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
    const ref = useRef();
    const src = `${sysConfig.webServer}/api/v1/image?hash=${hash}`;

    const [blob, setBlob] = useState();
    const [url, setUrl] = useState();
    const controller = useRef(new AbortController());
    useEffect(() => {
        fetch(src, {
            method: 'get',
            signal: controller.current.signal,
        }).then(response => {
            return response.blob();
        }).then(b => {
            setBlob(b);
            const objectURL = URL.createObjectURL(b);
            console.log(b, objectURL);
            setUrl(objectURL);
        });
    }, []);

    const shareImage = async () => {
        try {
            console.log('shareImage', blob);
            if (Platform.OS === "web") {
                const b = new Blob([blob], { type: "image/png" });
                console.log(b);
                navigator.clipboard.write([
                    new ClipboardItem({
                        [b.type]: b
                    })
                ]);
                window.noteInfo(`拷贝图片成功：${hash}`);
                return
            }
            await Sharing.shareAsync({ url: url });
        } catch (e) {
            console.log(e);
        }
    };

    const downloadImage = async () => {
        try {
            console.log('downloadImage', blob);
            if (Platform.OS === "web") {
                const b = new Blob([blob], { type: "image/png" });
                downloadURI(url, `${hash}.png`);
                window.URL.revokeObjectURL(url);
                window.noteInfo(`下载图片成功：${hash}`);
                return
            }
            await Sharing.shareAsync({ url: url });
        } catch (e) {
            console.log(e);
        }
    };
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
                <Appbar.Action icon="export-variant" onPress={shareImage} />
                <Appbar.Action icon="download" onPress={downloadImage} />
            </Appbar.Header>
            <ViewShot ref={ref} style={{
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
                }} source={url} />
            </ViewShot>
        </>
    );
}

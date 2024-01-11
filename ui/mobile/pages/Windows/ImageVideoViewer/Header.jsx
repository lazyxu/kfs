import { listDriverFileByHash } from "@kfs/common/api/webServer/driverfile";
import * as MediaLibrary from 'expo-media-library';
import { useEffect, useState } from "react";
import { View } from 'react-native';
import { Badge, IconButton, Surface, Text } from "react-native-paper";

export default function ({ navigation, hash, uri, index, total }) {
    const [sameFiles, setSameFiles] = useState([]);
    const [downloadName, setDownloadName] = useState();
    useEffect(() => {
        listDriverFileByHash(hash).then(sf => {
            setSameFiles(sf);
            setDownloadName(sf[0].name);
        });
    }, [hash]);
    return (
        <Surface style={{
            flexDirection: 'row',
            alignItems: 'center',
            justifyContent: "space-between",
            paddingHorizontal: 4,
            position: "absolute", left: 0, top: 0, right: 0,
            zIndex: 1,
        }}>
            <IconButton
                style={{ borderRadius: 0, backgroundColor: null, transition: null }}
                onPress={() => navigation.pop()}
                // iconColor={actionIconColor}
                icon="chevron-left"
            // disabled={disabled}
            // rippleColor={rippleColor}
            />
            <Text>{index+1}/{total}</Text>
            <View style={{
                flexDirection: 'row',
                alignItems: 'center',
            }}>
                <View>
                    <IconButton
                        style={{ borderRadius: 0, backgroundColor: null, transition: null }}
                        onPress={() => navigation.navigate("SameFile", { hash, sameFiles })}
                        icon="pound"
                    />
                    <Badge style={{ position: 'absolute', top: 4, right: 0, }}>
                        {sameFiles.length}
                    </Badge>
                </View>
                <IconButton
                    style={{ borderRadius: 0, backgroundColor: null, transition: null }}
                    onPress={async () => {
                        const writeOnly = false;
                        let perm = await MediaLibrary.getPermissionsAsync(writeOnly);
                        console.log("getPermissionsAsync", perm);
                        if (!perm.granted) {
                            console.log("requestPermissionsAsync");
                            perm = await MediaLibrary.requestPermissionsAsync(writeOnly);
                            console.log("requestPermissionsAsync", perm);
                        }
                        if (!perm.granted) {
                            window.noteWarning("保存失败，无权限");
                            return;
                        }
                        const asset = await MediaLibrary.createAssetAsync(uri);
                        const albumName = "考拉云盘";
                        let album = await MediaLibrary.getAlbumAsync(albumName);
                        if (album) {
                            await MediaLibrary.addAssetsToAlbumAsync(asset, album, true);
                        } else {
                            await MediaLibrary.createAlbumAsync(albumName, asset, true);
                        }
                        window.noteInfo("保存成功");
                    }}
                    icon="download-outline"
                    disabled={!downloadName}
                />
                <IconButton
                    // size={size}
                    // onPress={downloadImage}
                    // iconColor={actionIconColor}
                    icon="dots-vertical"
                // disabled={disabled}
                // rippleColor={rippleColor}
                />
            </View>
        </Surface>
    );
}

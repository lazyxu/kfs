import { listDriverFileByHash } from "@kfs/common/api/webServer/driverfile";
import { useEffect, useState } from "react";
import { View } from 'react-native';
import { Badge, IconButton, Surface, Text } from "react-native-paper";

export default function ({ navigation, hash, index, total }) {
    const [sameFiles, setSameFiles] = useState([]);
    const [downloadName, setDownloadName] = useState();
    useEffect(() => {
        listDriverFileByHash(hash).then(sf => {
            setSameFiles(sf);
            setDownloadName(sf[0].name);
        });
    }, []);
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
            <Text>{index}/{total}</Text>
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
                    // size={size}
                    // onPress={downloadImage}
                    // iconColor={actionIconColor}
                    icon="download-outline"
                // disabled={disabled}
                // rippleColor={rippleColor}
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

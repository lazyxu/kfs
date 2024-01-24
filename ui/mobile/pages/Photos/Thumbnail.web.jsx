import ImgCancelable from '@kfs/common/components/ImgCancelable';
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { Image, Pressable, View } from 'react-native';
import { Text } from 'react-native-paper';

function formatDuration(seconds) {
    seconds = Math.floor(seconds);
    const s = seconds % 60;
    const m = Math.floor(seconds / 60) % 60;
    const h = Math.floor(seconds / 3600) % 60;
    if (h > 0) {
        return "" + h + ":" + (m < 10 ? "0" + m : m) + ":" + (s < 10 ? "0" + s : s);
    }
    if (m > 0) {
        return "" + m + ":" + (s < 10 ? "0" + s : s);
    }
    return "0:" + (s < 10 ? "0" + s : s);
}

export default function ({ width, navigation, list, index, onLoaded }) {
    const src = `${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${list[index].hash}`;
    return (
        <Pressable onPress={() => navigation.navigate("Viewer", { list, index })}>
            <ImgCancelable
                src={src} onLoaded={onLoaded}
                renderImg={(url) => <Image style={{
                    width: width,
                    height: width,
                }} source={{ uri: url }} />}
                renderSkeleton={() => <View style={{
                    width: width,
                    height: width,
                }} />}
            />
            {list[index].type === "video" &&
                <View style={{
                    position: "absolute",
                    height: "100%",
                    width: "100%",
                    alignItems: 'flex-end',
                    justifyContent: 'end',
                }}>
                    <Text style={{ color: "white" }}>
                        {formatDuration(list[index].duration)}
                    </Text>
                </View>
            }
        </Pressable>
    );
}

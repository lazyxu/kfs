import { Pressable } from 'react-native';
import { Appbar, Surface, Text } from "react-native-paper";

export default function ({ navigation, route }) {
    const { hash, sameFiles } = route.params;
    return (
        <>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="相同文件" />
                <Pressable onPress={() => navigation.pop()}  >
                    <Text>完成</Text>
                </Pressable>
            </Appbar.Header>
            <Surface style={{
                flex: 1,
                whiteSpace: "pre"
            }}>
                <Text>共 {sameFiles.length} 个相同文件</Text>
                {sameFiles.map((f, i) => <Text key={i}>
                    {f.driverName}:{f.dirPath.length ? ("/" + f.dirPath.join("/") + "/" + f.name) : ("/" + f.name)}
                </Text>)}
            </Surface>
        </>
    );
}

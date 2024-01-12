import { Appbar, List, Surface } from 'react-native-paper';

export default function () {
    const navigation = window.kfsNavigation;
    return (
        <Surface>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="搜索" />
            </Appbar.Header>
            <List.Section>
                <List.Accordion style={{ padding: 0, margin: 0 }}  title="位置" />
                <List.Accordion style={{ padding: 0, margin: 0 }}  title="文件类型"
                    left={props => <List.Icon {...props} icon="file-image" />}
                    right={() => <></>} onPress={() => navigation.navigate("SearchType")}
                />
                <List.Accordion style={{ padding: 0, margin: 0 }}  title="人物识别" />
                <List.Accordion style={{ padding: 0, margin: 0 }}  title="文本识别" />
                <List.Accordion style={{ padding: 0, margin: 0 }}  title="物体识别" />

                <List.Accordion style={{ padding: 0, margin: 0 }}  title="拍摄设备" />
                <List.Accordion style={{ padding: 0, margin: 0 }}  title="文件大小" />
                <List.Accordion style={{ padding: 0, margin: 0 }}  title="文件后缀"
                    left={props => <List.Icon {...props} icon="file-jpg-box" />}
                    right={() => <></>} onPress={() => navigation.navigate("SearchSuffix")}
                />
            </List.Section>
        </Surface>
    );
}

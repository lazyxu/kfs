import { Appbar, List } from 'react-native-paper';
import Suffix from "./Suffix";
import Type from './Type';

export default function () {
    return (
        <>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="搜索" />
            </Appbar.Header>
            <List.Section>
                <List.Accordion title="位置" />
                <Type />
                <List.Accordion title="人物识别" />
                <List.Accordion title="文本识别" />
                <List.Accordion title="物体识别" />

                <List.Accordion title="拍摄设备" />
                <List.Accordion title="文件大小" />
                <Suffix />
            </List.Section>
        </>
    );
}

import { listDCIMMetadataTime } from '@kfs/common/api/webServer/exif';
import { Appbar, Surface } from "react-native-paper";
import ThumbnailListAll from './ThumbnailListAll';

export default function () {
    return (
        <Surface style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="照片" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <ThumbnailListAll listDCIMMetadataTime={listDCIMMetadataTime} />
        </Surface>
    );
}

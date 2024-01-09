import { listDCIMMetadataTime } from '@kfs/common/api/webServer/exif';
import { Appbar } from "react-native-paper";
import ThumbnailList from './ThumbnailList';

export default function () {
    return (
        <>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="照片" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <ThumbnailList listDCIMMetadataTime={listDCIMMetadataTime}/>
        </>
    );
}

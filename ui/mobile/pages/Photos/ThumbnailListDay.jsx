import ThumbnailListTemplate from "./ThumbnailListTemplate";

const elementsPerLine = 3;
function getTag(m) {
    return m.year + " 年 " + m.month + " 月 " + m.day + " 日";
}

export default function ({ metadataList, listDCIMMetadataTime }) {
    return (
        <ThumbnailListTemplate metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} getTag={getTag} elementsPerLine={elementsPerLine} />
    );
}

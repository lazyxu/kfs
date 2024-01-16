import ThumbnailListTemplate from "./ThumbnailListTemplate";

const elementsPerLine = 5;
function getTag(m) {
    return m.year + " 年 " + m.month + " 月";
}

export default function ({ metadataList, listDCIMMetadataTime }) {
    return (
        <ThumbnailListTemplate metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} getTag={getTag} elementsPerLine={elementsPerLine} />
    );
}

import ThumbnailListTemplate from "./ThumbnailListTemplate";

const elementsPerLine = 10;
function getTag(m) {
    return m.year + " 年 ";
}

export default function ({ metadataList, listDCIMMetadataTime }) {
    return (
        <ThumbnailListTemplate metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} getTag={getTag} elementsPerLine={elementsPerLine} />
    );
}

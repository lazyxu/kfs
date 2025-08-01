import ThumbnailListTemplate from "./ThumbnailListTemplate";

const elementsPerLine = 5;
function getTag(m) {
    return m.year + " 年 " + m.month + " 月";
}

export default function ({ metadataList }) {
    return (
        <ThumbnailListTemplate metadataList={metadataList} getTag={getTag} elementsPerLine={elementsPerLine} />
    );
}

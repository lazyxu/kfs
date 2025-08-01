import ThumbnailListTemplate from "./ThumbnailListTemplate";

const elementsPerLine = 10;
function getTag(m) {
    return m.year + " 年";
}

export default function ({ metadataList }) {
    return (
        <ThumbnailListTemplate metadataList={metadataList} getTag={getTag} elementsPerLine={elementsPerLine} />
    );
}

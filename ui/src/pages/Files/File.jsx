import AbsolutePath from "components/AbsolutePath";
import FileViewer from "./FileViewer/FileViewer";
import Dialog from "components/Dialog";

export default function ({file}) {
    return (
        <>
            <AbsolutePath/>
            <FileViewer file={file}/>
            <Dialog/>
        </>
    );
}

import React from 'react';
import NewFile from "./NewFile";
import useDialog from "hox/dialog";
import NewDir from "./NewDir";
import DialogAttribute from "./DialogAttribute";
import DialogNewDriver from "./DialogNewDriver";

export default function () {
    const [dialog, setDialog] = useDialog();
    console.log("dialog", dialog);
    // if (dialog === null || dialog === undefined) {
    //     return <div/>
    // }
    // return (dialog?.title === "新建文件" && <NewFile/>)
    switch (dialog?.title) {
        case "新建文件":
            return <NewFile />;
        case "新建文件夹":
            return <NewDir />;
        case "新建云盘":
            return <DialogNewDriver />;
        case "属性":
            return <DialogAttribute />;
    }
    return <div />;
};

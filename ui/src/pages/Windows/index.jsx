import useWindows, { APP_TEXT_VIEWER } from "hox/windows";
import TextViewer from "./TextViewer";

export default function () {
    const [windows, setWindows] = useWindows();
    return (
        Object.values(windows).map(w => (
            <div key={w.id}>
                {w.app === APP_TEXT_VIEWER && <TextViewer id={w.id} props={w.props} />}
            </div>
        ))
    )
}

import useResourceManager from 'hox/resourceManager';
import {list, openFile} from "api/fs";
import {Link, Stack} from "@mui/material";
import SvgIcon from "../../Icon/SvgIcon";
import {listDriver} from "api/driver";

export default ({type, name, icon, filePath}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Link color="inherit" underline="hover" onClick={(() => {
            if (filePath) {
                if (type === "file") {
                    openFile(setResourceManager, resourceManager.driverName, filePath, resourceManager.file);
                } else {
                    list(setResourceManager, resourceManager.driverName, filePath);
                }
            } else {
                listDriver(setResourceManager);
            }
        })}>
            <Stack direction="row"
                   justifyContent="flex-start"
                   alignItems="center"
            >
                {icon && <SvgIcon icon={icon}/>}
                {name}
            </Stack>
        </Link>
    )
};

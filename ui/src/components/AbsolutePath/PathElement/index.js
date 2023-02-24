import useResourceManager from 'hox/resourceManager';
import {list, openFile} from "api/fs";
import {Link, Stack} from "@mui/material";
import SvgIcon from "../../Icon/SvgIcon";
import {listBranch} from "api/branch";

export default ({type, name, icon, filePath}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Link color="inherit" underline="hover" onClick={(() => {
            if (filePath) {
                if (type === "file") {
                    openFile(setResourceManager, resourceManager.branchName, filePath, resourceManager.file);
                } else {
                    list(setResourceManager, resourceManager.branchName, filePath);
                }
            } else {
                listBranch(setResourceManager);
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

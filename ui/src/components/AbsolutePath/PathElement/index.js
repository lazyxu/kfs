import useResourceManager from 'hox/resourceManager';
import {useClick} from "use";
import {open} from "api/api";
import {Link, Stack} from "@mui/material";
import SvgIcon from "../../Icon/SvgIcon";

export default ({name, icon, filePath}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Link color="inherit" underline="hover" onMouseDown={useClick(() => {
            open(setResourceManager, resourceManager.branchName, filePath);
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

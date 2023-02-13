import useResourceManager from 'hox/resourceManager';
import {useClick} from "use";
import {open} from "api/api";
import useSysConfig from "hox/sysConfig";
import {Link, Stack} from "@mui/material";
import SvgIcon from '@mui/material/SvgIcon';

export default ({name, icon, filePath}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    return (
        <Link color="inherit" underline="hover" onMouseDown={useClick(() => {
            open(sysConfig, setResourceManager, resourceManager.branchName, filePath);
        })}>
            <Stack direction="row"
                   justifyContent="flex-start"
                   alignItems="center"
            >
                {icon && <SvgIcon color="inherit" fontSize="inherit">
                    <svg
                        aria-hidden="true"
                        viewBox="0 0 200 200"
                        preserveAspectRatio="xMinYMin meet"
                    >
                        <use xlinkHref={`#icon-${icon}`}/>
                    </svg>
                </SvgIcon>}
                {name}
            </Stack>
        </Link>
    )
};

import SvgIcon from "@mui/material/SvgIcon";

export default function ({icon, color, className}) {
    return (
        <SvgIcon color={color ? color : "inherit"} fontSize="inherit" className={className}>
            <svg
                aria-hidden="true"
                viewBox="0 0 200 200"
                preserveAspectRatio="xMinYMin meet"
            >
                <use xlinkHref={`#icon-${icon}`}/>
            </svg>
        </SvgIcon>
    );
}

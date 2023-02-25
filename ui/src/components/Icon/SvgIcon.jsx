import SvgIcon from "@mui/material/SvgIcon";

export default function ({icon, color, className, style}) {
    return (
        <SvgIcon color={color ? color : "inherit"} fontSize="inherit" className={className} style={style}>
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

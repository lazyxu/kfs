import useResourceManager, { openDir, openDrivers } from '@kfs/common/hox/resourceManager';
import SvgIcon from "@kfs/mui/components/Icon/SvgIcon";
import { ArrowBack, ArrowForward, ArrowForwardIos, ArrowUpward, Refresh, Search } from "@mui/icons-material";
import { Button, IconButton, Menu, MenuItem, Stack, TextField } from "@mui/material";
import { useTheme } from "@mui/material/styles";
import { useRef, useState } from "react";

function PathButton({ name, icon, dirPath, onMenu }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver } = resourceManager;
    const [anchorEl, setAnchorEl] = useState(null);

    const handleContextMenu = (e) => {
        e.preventDefault(); e.stopPropagation();
        setAnchorEl(e.currentTarget);
    };

    const handleMenuClose = () => setAnchorEl(null);

    return (
        <>
            <Button
                variant="text"
                size="small"
                startIcon={icon ? <SvgIcon icon={icon} sx={{ verticalAlign: "middle", px: 1 }} /> : null}
                sx={{
                    borderRadius: 1,
                    minWidth: 0,
                    px: 1,
                    bgcolor: theme => theme.background.primary,
                    color: theme => theme.context.secondary,
                    display: "flex",
                    alignItems: "center",
                    gap: 0.5,
                    boxShadow: "none",
                    ":hover": { bgcolor: theme => theme.palette.action.hover }
                }}
                onClick={e => {
                    e.preventDefault(); e.stopPropagation();
                    if (Array.isArray(dirPath)) {
                        openDir(setResourceManager, driver, dirPath);
                    } else {
                        openDrivers(setResourceManager);
                    }
                }}
                onContextMenu={handleContextMenu}
            >
                <span style={{ lineHeight: "normal" }}>{name}</span>
            </Button>
            <IconButton
                size="small"
                sx={{
                    borderRadius: 1,
                    minWidth: 0,
                    bgcolor: theme => theme.background.primary,
                    color: theme => theme.context.secondary,
                    boxShadow: "none",
                    borderLeft: "none",
                    borderRight: theme => `1px solid ${theme.palette.divider}`,
                    m: 0,
                    width: 32,
                    height: 32,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    ":hover": { bgcolor: theme => theme.palette.action.hover }
                }}
                onClick={e => {
                    e.preventDefault(); e.stopPropagation();
                    setAnchorEl(e.currentTarget);
                }}
            >
                <ArrowForwardIos fontSize="small" />
            </IconButton>
            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
            >
                <MenuItem onClick={() => { e.preventDefault(); e.stopPropagation(); handleMenuClose(); onMenu?.("newFolder", dirPath); }}>新建文件夹</MenuItem>
                <MenuItem onClick={() => { e.preventDefault(); e.stopPropagation(); handleMenuClose(); onMenu?.("paste", dirPath); }}>粘贴</MenuItem>
                {/* 可扩展更多菜单项 */}
            </Menu>
        </>
    );
}

export default () => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const dirPath = resourceManager.dirPath || [];
    const [editMode, setEditMode] = useState(false);
    const [inputValue, setInputValue] = useState("");
    const inputRef = useRef();
    const theme = useTheme();
    const [searchMode, setSearchMode] = useState(false);
    const [searchValue, setSearchValue] = useState("");
    // 历史记录
    const [history, setHistory] = useState([dirPath]);
    const [historyIndex, setHistoryIndex] = useState(0);

    // 跳转目录并维护历史
    const jumpDir = (newPath) => {
        openDir(setResourceManager, resourceManager.driver, newPath);
        setHistory(prev => {
            const updated = prev.slice(0, historyIndex + 1).concat([newPath]);
            setHistoryIndex(updated.length - 1);
            return updated;
        });
    };

    // 后退
    const handleBack = () => {
        if (historyIndex > 0) {
            setHistoryIndex(historyIndex - 1);
            openDir(setResourceManager, resourceManager.driver, history[historyIndex - 1]);
        }
    };

    // 前进
    const handleForward = () => {
        if (historyIndex < history.length - 1) {
            setHistoryIndex(historyIndex + 1);
            openDir(setResourceManager, resourceManager.driver, history[historyIndex + 1]);
        }
    };

    // 上级
    const handleUp = () => {
        if (dirPath.length > 0) {
            const upPath = dirPath.slice(0, -1);
            jumpDir(upPath);
        }
    };

    // 刷新
    const handleRefresh = () => {
        openDir(setResourceManager, resourceManager.driver, dirPath, { refresh: true });
    };

    // 搜索
    const handleSearch = () => {
        // TODO: 实现搜索逻辑
        setSearchMode(false);
        // openDir/setResourceManager等，带搜索参数
    };

    // 编辑路径
    const handleEdit = () => {
        setInputValue(dirPath.join("\\"));
        setEditMode(true);
        setTimeout(() => inputRef.current?.focus(), 100);
    };
    const handleInputChange = (e) => setInputValue(e.target.value);
    const handleInputBlur = () => setEditMode(false);
    const handleInputKeyDown = (e) => {
        if (e.key === "Enter") {
            const parts = inputValue.split("/");
            console.log("Navigating to:", parts);
            if (parts[0]?.length > 0) {
                openDir(setResourceManager, parts[0], parts.slice(1));
            } else {
                openDrivers(setResourceManager);
            }
            setEditMode(false);
        }
    };

    const handleMenuAction = (action, path) => {
        // TODO: 实现新建文件夹、粘贴等功能
    };

    // Win11风格分隔符
    const Separator = () => (
        <ArrowForwardIos />
    );

    return (
        <Stack
            sx={{
                padding: "0 1em",
                margin: "0.3em",
                cursor: editMode ? "text" : "pointer",
                bgcolor: theme.palette.background.paper,
                borderRadius: 2,
                boxShadow: 1,
                minHeight: 48,
            }}
            direction="row"
            alignItems="center"
            spacing={0.5}
        >
            {/* Win11风格左侧操作按钮区 */}
            <Stack direction="row" alignItems="center" spacing={0.5} sx={{ minWidth: 160 }}>
                <IconButton size="small" onClick={handleBack} disabled={historyIndex === 0} title="后退">
                    <ArrowBack />
                </IconButton>
                <IconButton size="small" onClick={handleForward} disabled={historyIndex === history.length - 1} title="前进">
                    <ArrowForward />
                </IconButton>
                <IconButton size="small" onClick={handleUp} disabled={dirPath.length === 0} title="上级">
                    <ArrowUpward />
                </IconButton>
                <IconButton size="small" onClick={handleRefresh} title="刷新">
                    <Refresh />
                </IconButton>
            </Stack>
            {/* 中间地址栏 */}
            <Stack
                direction="row"
                alignItems="center"
                spacing={0.5}
                sx={{ flex: 1, minWidth: 0, position: "relative" }}
                onClick={e => {
                    if (!editMode) {
                        setInputValue((resourceManager.driver?.name || "") + "/" + dirPath.join("/"));
                        setEditMode(true);
                        setTimeout(() => inputRef.current?.focus(), 100);
                    }
                }}
            >
                {editMode ? (
                    <TextField
                        inputRef={inputRef}
                        value={inputValue}
                        onChange={handleInputChange}
                        onBlur={handleInputBlur}
                        onKeyDown={handleInputKeyDown}
                        size="small"
                        sx={{
                            width: "100%",
                            bgcolor: theme.palette.action.hover,
                            borderRadius: 2,
                            flex: 1,
                        }}
                        autoFocus
                    />
                ) : (
                    <Stack
                        direction="row"
                        alignItems="center"
                        spacing={0.5}
                        sx={{
                            flex: 1,
                            minWidth: 0,
                            bgcolor: theme.palette.action.hover,
                            borderRadius: 2,
                            border: `1px solid ${theme.palette.divider}`,
                            px: 1,
                            py: 0.5,
                        }}
                    >
                        <PathButton name="我的云盘" icon="wangpan" dirPath={[]} isCurrent={false} onMenu={handleMenuAction} />
                        {resourceManager.driver?.name && <PathButton name={resourceManager.driver?.name} icon="folder" dirPath={[]} isCurrent={dirPath.length === 0} onMenu={handleMenuAction} />}
                        {dirPath.map((name, i) => (
                            <Stack direction="row" alignItems="center" key={i}>
                                <PathButton
                                    name={name}
                                    dirPath={dirPath.slice(0, i + 1)}
                                    isCurrent={i === dirPath.length - 1}
                                    onMenu={handleMenuAction}
                                />
                            </Stack>
                        ))}
                    </Stack>
                )}
            </Stack>
            {/* 右侧搜索栏，始终显示搜索按钮，点击后显示输入框 */}
            <Stack direction="row" alignItems="center" spacing={0.5} sx={{ minWidth: 200, flex: "0 0 320px" }}>
                <TextField
                    value={searchValue}
                    onChange={e => setSearchValue(e.target.value)}
                    onKeyDown={e => {
                        if (e.key === "Enter") handleSearch();
                    }}
                    size="small"
                    sx={{
                        minWidth: 160,
                        maxWidth: 320,
                        bgcolor: theme.palette.action.hover,
                        borderRadius: 2,
                        flex: 1,
                        mr: 1,
                    }}
                    placeholder="搜索"
                    InputProps={{
                        endAdornment: (
                            <IconButton
                                size="small"
                                onClick={handleSearch}
                                title="搜索"
                                sx={{
                                    p: 0.5,
                                    bgcolor: theme.palette.action.hover,
                                    borderRadius: 1,
                                    "&:hover": { bgcolor: theme.palette.action.selected }
                                }}
                            >
                                <Search />
                            </IconButton>
                        )
                    }}
                />
            </Stack>
        </Stack>
    );
};

import { Box, List, ListItem, ListItemButton, ListItemIcon, ListItemText } from "@mui/material"
import AutoDeleteIcon from '@mui/icons-material/AutoDelete';
import TocIcon from '@mui/icons-material/Toc';
import "./Sidebar.css"

export default function Sidebar() {
    return (
        <Box sx={{ minWidth: '150px', height: '100%', bgcolor: 'rgb(34, 117, 240)' }}>
            <nav aria-label="main mailbox folders">
                <List>
                <ListItem disablePadding>
                    <ListItemButton>
                        <ListItemIcon><TocIcon /></ListItemIcon>
                        <ListItemText primary="Table" />
                    </ListItemButton>
                </ListItem>
                <ListItem disablePadding >
                    <ListItemButton>
                        <ListItemIcon><AutoDeleteIcon /></ListItemIcon>
                        <ListItemText primary="Delete" />
                    </ListItemButton>
                </ListItem>
                </List>
            </nav>
        </Box>
    )
}
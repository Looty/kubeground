import "./Table.css"
import CustomTable from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';

function createData(
  level: string,
  resources: number,
  isDeployed: string,
  isCompleted: string,
) {
  return { level, resources, isDeployed, isCompleted };
}

const rows = [
  createData('1', 2, 'false', 'true'),
  createData('2', 1, 'false', 'false'),
  createData('3', 5, 'true', 'false'),
];

export default function Table() {
    return (
        <TableContainer component={Paper}>
            <CustomTable sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
                <TableHead>
                <TableRow>
                    <TableCell>Level</TableCell>
                    <TableCell>Resources num</TableCell>
                    <TableCell>Deployed?</TableCell>
                    <TableCell>Completed?</TableCell>
                </TableRow>
                </TableHead>
                <TableBody>
                {rows.map((row) => (
                    <TableRow key={row.level} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                        <TableCell component="th" scope="row">
                            {row.level}
                        </TableCell>
                        <TableCell>{row.resources}</TableCell>

                        <TableCell><CheckIcon /></TableCell>
                        <TableCell><CloseIcon /></TableCell>
                    </TableRow>
                ))}
                </TableBody>
            </CustomTable>
        </TableContainer>
    )
}
import { Box, Grid } from '@mui/material'
import './App.css'
import Footer from './layout/footer/Footer'
import Main from './layout/main/Main'
import Navbar from './layout/navbar/Navbar'
import Sidebar from './layout/sidebar/Sidebar'
import Table from './pages/table/Table'

function App() {

  return (
   <>
   <Navbar />
   <Box sx={{ flexGrow: 1 }}>
    <Grid container spacing={0}>
      <Grid item xs={1}>
        <Sidebar />
      </Grid>
      <Grid item xs={11}>
        <Main />
        {/* <Table /> */}
      </Grid>
    </Grid>
   </Box>
   <Footer />
   </>
  )
}

export default App

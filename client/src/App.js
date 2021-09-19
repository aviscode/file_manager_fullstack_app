import React from 'react';
import DeleteIcon from '@material-ui/icons/Delete';
import ArrowDownwardIcon from '@material-ui/icons/ArrowDownward';
import {
    Box,
    Button,
    IconButton,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography
} from "@material-ui/core";

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            items: [],
            isLoaded: false,
            selectedFile: null,
            fileName: "",
        };
    }

    componentDidMount() {
        this.getFiles()
    }

    getFiles = () => {
        fetch('files')
            .then(res => res.json())
            .then(result => {
                this.setState({
                    isLoaded: true,
                    items: result
                });
            });
    }

    onFileUpload = () => {
        if (this.state.selectedFile === null) {
            return
        }
        const formData = new FormData();
        formData.append("file", this.state.selectedFile);
        fetch("files", {method: 'PUT', body: formData,})
            .then(this.getFiles, this.setState({selectedFile: null}))
            .catch((error) => {
                console.error('Error:', error);
            });
    };

    onFileChange = event => {
        // update the state
        this.setState({selectedFile: event.target.files[0]});
        this.setState({filename: event.target.files[0] ? event.target.files[0].name : ""});
    };

    convertBytes = function (bytes) {
        const sizes = ["Bytes", "KB", "MB", "GB", "TB"]

        if (bytes === 0) {
            return "n/a"
        }

        const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)))

        if (i === 0) {
            return bytes + " " + sizes[i]
        }

        return (bytes / Math.pow(1024, i)).toFixed(1) + " " + sizes[i]
    }


    render() {
        const {items} = this.state;
        return (
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                minHeight="100vh">
                <TableContainer component={Paper} style={{maxWidth: '1000px', minWidth: '550px'}}>
                    <Box display="flex"
                         justifyContent="center"
                         alignItems="center">
                        <Button variant="contained" color="primary" component="label">
                            Choose a File
                            <input type="file" hidden onChange={this.onFileChange}/>
                        </Button>
                        <Typography variant="h6" component="h2" hidden={this.state.selectedFile === null}>
                            &nbsp;&nbsp;&nbsp; file: &nbsp; {this.state.filename} &nbsp;&nbsp;&nbsp;
                            <Button
                                color="primary"
                                variant="contained"
                                component="span"
                                onClick={this.onFileUpload}>
                                Upload
                            </Button>
                        </Typography>


                    </Box>
                    <Table aria-label="simple table">
                        <TableHead>
                            <TableRow>
                                <TableCell>File Name</TableCell>
                                <TableCell align="right">Size</TableCell>
                                <TableCell align="right">Format</TableCell>
                                <TableCell align="right">Uploaded</TableCell>
                                <TableCell align="center">Operations</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {!items ? null : items.map((row) => (
                                <TableRow key={row.filename}>
                                    <TableCell component="th" scope="row">
                                        {row.filename}
                                    </TableCell>
                                    <TableCell align="right">{this.convertBytes(row.size)}</TableCell>
                                    <TableCell align="right">{row.format}</TableCell>
                                    <TableCell align="right">{row.uploaded}</TableCell>
                                    <TableCell align="center">{
                                        <div>
                                            <IconButton aria-label="get"
                                                        onClick={() => fetch("files/" + row.filename)
                                                            .then(res => {
                                                                return res.blob();
                                                            }).then(blob => {
                                                                const href = window.URL.createObjectURL(blob);
                                                                const link = document.createElement('a');
                                                                link.href = href;
                                                                link.setAttribute('download', row.filename); //or any other extension
                                                                document.body.appendChild(link);
                                                                link.click();
                                                                document.body.removeChild(link);
                                                            }).catch(err => console.error(err))
                                                        } color="primary">
                                                <ArrowDownwardIcon/>
                                            </IconButton>
                                            <IconButton aria-label="delete" color="secondary"
                                                        onClick={() => fetch("files/" + row.filename, {method: 'DELETE'}).then(this.getFiles)}>
                                                <DeleteIcon/>
                                            </IconButton>
                                        </div>
                                    }
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Box>
        );
    }
}

export default App;
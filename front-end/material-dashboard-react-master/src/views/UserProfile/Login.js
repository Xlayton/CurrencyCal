import React, { useState } from "react";
// @material-ui/core components
import { makeStyles } from "@material-ui/core/styles";
import Input from "@material-ui/core/Input";
// core components
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import Button from "components/CustomButtons/Button.js";
import Card from "components/Card/Card.js";
import CardHeader from "components/Card/CardHeader.js";
import CardBody from "components/Card/CardBody.js";
import CardFooter from "components/Card/CardFooter.js";

const styles = {
    cardCategoryWhite: {
        color: "rgba(255,255,255,.62)",
        margin: "0",
        fontSize: "14px",
        marginTop: "0",
        marginBottom: "0"
    },
    cardTitleWhite: {
        color: "#FFFFFF",
        marginTop: "0px",
        minHeight: "auto",
        fontWeight: "300",
        fontFamily: "'Roboto', 'Helvetica', 'Arial', sans-serif",
        marginBottom: "3px",
        textDecoration: "none"
    }
};

const useStyles = makeStyles(styles);

export default function LoginPage() {
    const classes = useStyles();

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [user_id, setUserId] = useState("");

    return (
        <div>
            <GridContainer>
                <GridItem xs={12} sm={12} md={12}>
                    <Card>
                        <CardHeader color="success">
                            <h4 className={classes.cardTitleWhite}>Login</h4>
                        </CardHeader>
                        <CardBody>
                            <GridContainer>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Cardholder Id (4 digits)"
                                        id="username"
                                        value={username}
                                        onChange={(evt) => { setUsername(evt.target.value); console.log(username) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Password"
                                        id="password"
                                        value={password}
                                        onChange={(evt) => { setPassword(evt.target.value); console.log(password) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <h4>Your unique Id. Please copy for the user profile page: {user_id}</h4>
                                </GridItem>
                            </GridContainer>
                        </CardBody>
                        <CardFooter>
                            <Button color="success" onClick={() => {
                                fetch("http://35.188.19.111/assistlogin", {
                                    method: "POST",
                                    headers: {
                                        "Content-Type": "application/json"
                                    },
                                    body: JSON.stringify({
                                        "cardholder_id": parseInt(username, 10)
                                    })
                                })
                                .then(res => res.json())
                                .then(data => setUserId(data.user_uuid))
                            }}>Login</Button>
                        </CardFooter>
                    </Card>
                </GridItem>
            </GridContainer>
        </div>
    );
}

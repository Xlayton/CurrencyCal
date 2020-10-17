import React, { useState } from "react";
// @material-ui/core components
import { makeStyles } from "@material-ui/core/styles";
import InputLabel from "@material-ui/core/InputLabel";
import BugReport from "@material-ui/icons/BugReport";
import Code from "@material-ui/icons/Code";
import Cloud from "@material-ui/icons/Cloud";

// core components
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import CustomInput from "components/CustomInput/CustomInput.js";
import Input from "@material-ui/core/Input";
import Button from "components/CustomButtons/Button.js";
import Card from "components/Card/Card.js";
import CardHeader from "components/Card/CardHeader.js";
import CardAvatar from "components/Card/CardAvatar.js";

// @material-ui/core components
import Table from "components/Table/Table.js";
import Tasks from "components/Tasks/Tasks.js";
import CustomTabs from "components/CustomTabs/CustomTabs.js";
import Danger from "components/Typography/Danger.js";
import CardIcon from "components/Card/CardIcon.js";
import CardBody from "components/Card/CardBody.js";
import CardFooter from "components/Card/CardFooter.js"
// @material-ui/icons
import Edit from "@material-ui/icons/Edit";
import Close from "@material-ui/icons/Close";
import Check from "@material-ui/icons/Check";

import { bugs, website, server, agreementTask, agreementUrl } from "variables/general.js";
import avatar from "assets/img/faces/marc.jpg";

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

export default function RegisterProfile() {
    const classes = useStyles();
    const [address, setAddress] = useState();
    const [city, setCity] = useState();
    const [state, setState] = useState();
    const [street, setStreet] = useState();
    const [unit, setUnit] = useState();
    const [zip_code, setZipCode] = useState();
    const [agreements, setAgreements] = useState();
    const [email, setEmail] = useState();
    const [first_name, setFirstName] = useState();
    const [identification, setIdentification] = useState();
    const [date_of_birth, setDateOfBirth] = useState();
    const [id, setId] = useState();
    const [id_type, setIdType] = useState();
    const [income, setIncome] = useState();
    const [amount, setAmount] = useState();
    const [frequency, setFrequency] = useState();
    const [occupation, setOccupation] = useState();
    const [source, setSource] = useState();
    const [last_name, setLastName] = useState();
    const [mobile, setMobile] = useState();
    const [shipping_address, setShippingAdress] = useState();
    const [shipping_city, setShippingCity] = useState();
    const [shipping_state, setShippingState] = useState();
    const [shipping_street, setShippingStreet] = useState();
    const [shipping_unit, setShippingUnit] = useState();
    const [shipping_zip_code, setShippingZipCode] = useState();

    return (
        <div>
            <GridContainer>
                <GridItem xs={12} sm={12} md={12}>
                    <Card>
                        <CardHeader color="primary">
                            <h4 className={classes.cardTitleWhite}>Edit Profile</h4>
                            <p className={classes.cardCategoryWhite}>Complete your profile</p>
                        </CardHeader>
                        <CardBody>
                            <GridContainer>
                                <GridItem xs={12} sm={12} md={5}>
                                    <Input
                                        placeholder="City"
                                        id="city"
                                        value={city}
                                        onChange={(evt) => { setCity(evt.target.value); console.log(city) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={3}>
                                    <Input
                                        placeholder="State"
                                        id="state"
                                        value={state}
                                        onChange={(evt) => { setState(evt.target.value); console.log(state) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Street"
                                        id="street"
                                        value={street}
                                        onChange={(evt) => { setStreet(evt.target.value); console.log(street) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={6}>
                                    <Input
                                        placeholder="Unit"
                                        id="unit"
                                        value={unit}
                                        onChange={(evt) => { setUnit(evt.target.value); console.log(unit) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={6}>
                                    <Input
                                        placeholder="zip code"
                                        id="zip_code"
                                        value={zip_code}
                                        onChange={(evt) => { setZipCode(evt.target.value); console.log(zip_code) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Email"
                                        id="email"
                                        value={email}
                                        onChange={(evt) => { setEmail(evt.target.value); console.log(email) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="First Name"
                                        id="first_name"
                                        value={first_name}
                                        onChange={(evt) => { setFirstName(evt.target.value); console.log(first_name) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Identification"
                                        id="identification"
                                        value={identification}
                                        onChange={(evt) => { setIdentification(evt.target.value); console.log(identification) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Date of birth"
                                        id="date_of_birth"
                                        value={date_of_birth}
                                        onChange={(evt) => { setDateOfBirth(evt.target.value); console.log(date_of_birth) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Id"
                                        id="id"
                                        value={id}
                                        onChange={(evt) => { setId(evt.target.value); console.log(id) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Id Type"
                                        id="id_type"
                                        value={id_type}
                                        onChange={(evt) => { setIdType(evt.target.value); console.log(id_type) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Income"
                                        id="income"
                                        value={income}
                                        onChange={(evt) => { setIncome(evt.target.value); console.log(income) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Amount"
                                        id="amount"
                                        value={amount}
                                        onChange={(evt) => { setAmount(evt.target.value); console.log(amount) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Frequency"
                                        id="frequency"
                                        value={frequency}
                                        onChange={(evt) => { setFrequency(evt.target.value); console.log(frequency) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Occupation"
                                        id="occupation"
                                        value={occupation}
                                        onChange={(evt) => { setOccupation(evt.target.value); console.log(occupation) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Source"
                                        id="source"
                                        value={source}
                                        onChange={(evt) => { setSource(evt.target.value); console.log(source) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Last Name"
                                        id="last_name"
                                        value={last_name}
                                        onChange={(evt) => { setLastName(evt.target.value); console.log(last_name) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Mobile"
                                        id="mobile"
                                        value={mobile}
                                        onChange={(evt) => { setMobile(evt.target.value); console.log(mobile) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Shipping City"
                                        id="shipping_city"
                                        value={shipping_city}
                                        onChange={(evt) => { setShippingCity(evt.target.value); console.log(shipping_city) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Shipping State"
                                        id="shipping_state"
                                        value={shipping_state}
                                        onChange={(evt) => { setShippingState(evt.target.value); console.log(shipping_state) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Shipping Street"
                                        id="shipping_street"
                                        value={shipping_street}
                                        onChange={(evt) => { setShippingStreet(evt.target.value); console.log(shipping_street) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Shipping Unit"
                                        id="shipping_unit"
                                        value={shipping_unit}
                                        onChange={(evt) => { setShippingUnit(evt.target.value); console.log(shipping_unit) }}
                                    />
                                </GridItem>
                                <GridItem xs={12} sm={12} md={4}>
                                    <Input
                                        placeholder="Shipping Zip Code"
                                        id="shipping_zip_code"
                                        value={shipping_zip_code}
                                        onChange={(evt) => { setShippingZipCode(evt.target.value); console.log(shipping_zip_code) }}
                                    />
                                </GridItem>
                            </GridContainer>
                            <GridContainer>
                                <GridItem xs={12} sm={12} md={12}>
                                    <CustomTabs
                                        title="Agree:"
                                        headerColor="primary"
                                        tabs={[
                                            {
                                                tabName: "Agreements",

                                                tabContent: (
                                                    <Tasks
                                                        checkedIndexes={[]}
                                                        tasksIndexes={[0, 1, 2]}
                                                        tasks={agreementTask}
                                                        link={agreementUrl}
                                                    />
                                                )
                                            }
                                        ]}
                                    />
                                </GridItem>
                            </GridContainer>
                            <GridContainer>
                                <GridItem xs={12} sm={12} md={12}>
                                    {/* <InputLabel style={{ color: "#AAAAAA" }}>About me</InputLabel> */}
                                </GridItem>
                            </GridContainer>
                        </CardBody>
                        <CardFooter>
                            <Button color="primary" onClick={() => { console.log('onClick'); }}>Update Profile</Button>
                        </CardFooter>
                    </Card>
                </GridItem>
            </GridContainer>
        </div>
    );
}

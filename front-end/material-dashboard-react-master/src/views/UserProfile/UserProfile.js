import React, { useState } from "react";
// @material-ui/core components
import { makeStyles } from "@material-ui/core/styles";
import InputLabel from "@material-ui/core/InputLabel";
import Input from "@material-ui/core/Input";

// react plugin for creating charts
import ChartistGraph from "react-chartist";
// @material-ui/core
import Icon from "@material-ui/core/Icon";
import CustomInput from "components/CustomInput/CustomInput.js";
import Button from "components/CustomButtons/Button.js";
import Typography from '@material-ui/core/Typography';

// @material-ui/icons
import Store from "@material-ui/icons/Store";
import Warning from "@material-ui/icons/Warning";
import DateRange from "@material-ui/icons/DateRange";
import LocalOffer from "@material-ui/icons/LocalOffer";
import Update from "@material-ui/icons/Update";
import ArrowUpward from "@material-ui/icons/ArrowUpward";
import AccessTime from "@material-ui/icons/AccessTime";
import Accessibility from "@material-ui/icons/Accessibility";
import BugReport from "@material-ui/icons/BugReport";
import Code from "@material-ui/icons/Code";
import Cloud from "@material-ui/icons/Cloud";

// core components
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import Table from "components/Table/Table.js";
import Tasks from "components/Tasks/Tasks.js";
import CustomTabs from "components/CustomTabs/CustomTabs.js";
import Danger from "components/Typography/Danger.js";
import Card from "components/Card/Card.js";
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from "components/Card/CardHeader.js";
import CardIcon from "components/Card/CardIcon.js";
import CardBody from "components/Card/CardBody.js";
import CardFooter from "components/Card/CardFooter.js";
import CardAvatar from "components/Card/CardAvatar.js";

import zane from "assets/img/faces/zane.jpg";
import josh from "assets/img/faces/josh.jpg";
import clayton from "assets/img/faces/clayton.png";
import eric from "assets/img/faces/eric.png";

import {
  dailySalesChart,
  emailsSubscriptionChart,
  completedTasksChart
} from "variables/charts.js";

import styles from "assets/jss/material-dashboard-react/views/dashboardStyle.js";
import { Assessment } from "@material-ui/icons";

const useStyles = makeStyles(styles);

export default function UserProfile() {
  const classes = useStyles();
  const [balance, setBalance] = useState(-1);
  const [user_info, setUserInfo] = useState({ contacts: [] });
  const [user_id, setUserId] = useState("");

  return (
    <div>
      {/* Balance */}
      <GridContainer>
        <GridItem xs={12} sm={12} md={4}>
          <Card>
            <CardHeader color="success" stats icon>
              <CardIcon color="success">
                <Icon>payment</Icon>
              </CardIcon>
              <p className={classes.cardCategory}>Current Balance</p>
              <h3 className={classes.cardTitle}>
                $ {balance}
              </h3>
            </CardHeader>
            <CardFooter stats>
              <div className={classes.stats}>
                <Update />
                Just Updated
              </div>
            </CardFooter>
          </Card>
        </GridItem>
        {/* Contacts */}
        <GridItem xs={12} sm={6} md={4}>
          <Card>
            <CardHeader color="info" stats icon>
              <CardIcon color="info">
                <Accessibility />
              </CardIcon>
              <p className={classes.cardCategory}>Contacts</p>
              <h3 className={classes.cardTitle}>
                {user_info.contacts.length}
              </h3>
            </CardHeader>
            <CardFooter stats>
              <div className={classes.stats}>
                <Icon>contact_page</Icon>
                Add more on Alexa!
              </div>
            </CardFooter>
          </Card>
          {/* Fill */}
        </GridItem>
        <GridItem xs={12} sm={6} md={4}>
          <Card>
            <CardHeader color="warning" stats icon>
              <CardIcon color="warning">
                <Icon>monetization_on</Icon>
              </CardIcon>
              <p className={classes.cardCategory}>Powered By</p>
              <h3 className={classes.cardTitle}>
                Galileo Instant
              </h3>
            </CardHeader>
            <CardFooter stats>
              <div className={classes.stats}>
                <Icon>assessment</Icon>
                API and Balance provided by Galileo
              </div>
            </CardFooter>
          </Card>
        </GridItem>
      </GridContainer>

      {/* This is the actual info of the page maybe? */}
      <GridContainer>
        <Card className={classes.root} variant="outlined">
          <CardContent>
            <Typography className={classes.title} color="textSecondary" gutterBottom>
              Please paste your id here
            </Typography>
            <GridItem xs={12} sm={12} md={12}>
              <Input
                placeholder="Id..."
                id="user_id"
                value={user_id}
                onChange={(evt) => { setUserId(evt.target.value); console.log(user_id) }}
              />
            </GridItem>
          </CardContent>
        </Card>
        <CardFooter>
          <Button color="success" onClick={() => {
            if (balance < 0) {
              fetch("http://35.188.19.111/getbalance", {
                method: "POST",
                headers: {
                  "Content-Type": "application/json"
                },
                body: JSON.stringify({
                  "account_uuid": user_id
                })
              })
                .then(res => res.json())
                .then(data => setBalance(data.balance))
            }

            if (!user_info.first_name) {
              fetch("http://35.188.19.111/accinfo", {
                method: "POST",
                headers: {
                  "Content-Type": "application/json"
                },
                body: JSON.stringify({
                  "account_uuid": user_id
                })
              })
                .then(res => res.json())
                .then(data => setUserInfo(data))
            }
          }}>Submit</Button>
        </CardFooter>
      </GridContainer>
    </div >
  );
}

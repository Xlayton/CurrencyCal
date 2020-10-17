// ##############################
// // // Tasks for TasksCard - see Dashboard view
// #############################

var bugs = [
  'Sign contract for "What are conference organizers afraid of?"',
  "Lines From Great Russian Literature? Or E-mails From My Boss?",
  "Flooded: One year later, assessing what was lost and what was found when a ravaging rain swept through metro Detroit",
  "Create 4 Invisible User Experiences you Never Knew About"
];
var website = [
  "Flooded: One year later, assessing what was lost and what was found when a ravaging rain swept through metro Detroit",
  'Sign contract for "What are conference organizers afraid of?"'
];
var server = [
  "Lines From Great Russian Literature? Or E-mails From My Boss?",
  "Flooded: One year later, assessing what was lost and what was found when a ravaging rain swept through metro Detroit",
  'Sign contract for "What are conference organizers afraid of?"'
];
var agreementTask = [
  `I agree to the "Electronic Communications" agreement page. 
  Link: "https://instant-assets.s3.amazonaws.com/electronic_communications_1_0_1.pdf"`,
  `I agree to the privacy policy from AmazonAWS and their terms of agreement. 
  Link: "https://instant-assets.s3.amazonaws.com/privacy_policy_1_0_1.pdf"`,
  `I agree to the deposit account agreement provided to me. 
  Link: "https://instant-assets.s3.amazonaws.com/deposit_account_agreement_1_0_1.pdf"`
];

var agreementUrl = "https://instant-assets.s3.amazonaws.com/privacy_policy_1_0_1.pdf";

module.exports = {
  // these 3 are used to create the tasks lists in TasksCard - Dashboard view
  bugs,
  website,
  server,
  agreementTask,
  agreementUrl
};

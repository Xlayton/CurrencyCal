# -*- coding: utf-8 -*-

# This is a High Low Guess game Alexa Skill.
# The skill serves as a simple sample on how to use the
# persistence attributes and persistence adapter features in the SDK.
import random
import logging
import os
import boto3

from ask_sdk_core.skill_builder import CustomSkillBuilder
from ask_sdk_core.utils import is_request_type, is_intent_name
from ask_sdk_core.handler_input import HandlerInput

from ask_sdk_model import Response
from ask_sdk_s3.adapter import S3Adapter

SKILL_NAME = 'Budget Buddy'
bucket_name = os.environ.get('S3_PERSISTENCE_BUCKET')
s3_client = boto3.client('s3',
                         region_name=os.environ.get('S3_PERSISTENCE_REGION'),
                         config=boto3.session.Config(signature_version='s3v4',s3={'addressing_style': 'path'}))
s3_adapter = S3Adapter(bucket_name=bucket_name, path_prefix="Media", s3_client=s3_client)
sb = CustomSkillBuilder(persistence_adapter=s3_adapter)

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)




@sb.request_handler(can_handle_func=is_request_type("LaunchRequest"))
def launch_request_handler(handler_input):
    #Handler for Skill Launch.
    # type: (HandlerInput) -> Response
    attr = handler_input.attributes_manager.persistent_attributes
    if not attr:
        attr['ended_session_count'] = 0
        attr['conversation_state'] = 'ENDED'

    handler_input.attributes_manager.session_attributes = attr

    speech_text = (
        "Hi im currency Kal. My job is to help you with your spending. What more do we want for description"
        "Would you like to talk about spending?")
    reprompt = "Say yes to continue talking to Kal abount spending"

    handler_input.response_builder.speak(speech_text).ask(reprompt)
    return handler_input.response_builder.response


@sb.request_handler(can_handle_func=is_intent_name("AMAZON.HelpIntent"))
def help_intent_handler(handler_input):
    """Handler for Help Intent."""
    # type: (HandlerInput) -> Response
    speech_text = (
       "Im currency Kal I can help you with your spending but first we need to create an account")
    reprompt = "Try saying login"

    handler_input.response_builder.speak(speech_text).ask(reprompt)
    return handler_input.response_builder.response


@sb.request_handler(
    can_handle_func=lambda input:
        is_intent_name("AMAZON.CancelIntent")(input) or
        is_intent_name("AMAZON.StopIntent")(input))
def cancel_and_stop_intent_handler(handler_input):
    """Single handler for Cancel and Stop Intent."""
    # type: (HandlerInput) -> Response
    speech_text = "Thanks for using Budget Buddy!!"

    handler_input.response_builder.speak(
        speech_text).set_should_end_session(True)
    return handler_input.response_builder.response


@sb.request_handler(can_handle_func=is_request_type("SessionEndedRequest"))
def session_ended_request_handler(handler_input):
    """Handler for Session End."""
    # type: (HandlerInput) -> Response
    logger.info(
        "Session ended with reason: {}".format(
            handler_input.request_envelope.request.reason))
    return handler_input.response_builder.response


def currently_playing(handler_input):
    """Function that acts as can handle for game state."""
    # type: (HandlerInput) -> bool
    is_currently_playing = False
    session_attr = handler_input.attributes_manager.session_attributes

    if ("conversation_state" in session_attr
            and session_attr['conversation_state'] == "STARTED"):
        is_currently_playing = True

    return is_currently_playing


@sb.request_handler(can_handle_func=lambda input:
                    not currently_playing(input) and
                    is_intent_name("AMAZON.YesIntent")(input))
def yes_handler(handler_input):
    """Handler for Yes Intent, only if the player said yes for
    a new game.
    """
    # type: (HandlerInput) -> Response
    session_attr = handler_input.attributes_manager.session_attributes
    session_attr['conversation_state'] = "STARTED"
    session_attr['guess_number'] = random.randint(0, 100)
    session_attr['no_of_guesses'] = 0

    speech_text = "Great! Try saying something to Kal if you arent already logged in"
    reprompt = "Try saying Register or login or if your already logged in try saying what is my balance"

    handler_input.response_builder.speak(speech_text).ask(reprompt)
    return handler_input.response_builder.response


@sb.request_handler(can_handle_func=lambda input:
                    not currently_playing(input) and
                    is_intent_name("AMAZON.NoIntent")(input))
def no_handler(handler_input):
    """Handler for No Intent, only if the player said no for
    a new game.
    """
    # type: (HandlerInput) -> Response
    session_attr = handler_input.attributes_manager.session_attributes
    session_attr['conversation_state'] = "ENDED"
    session_attr['ended_session_count'] += 1

    handler_input.attributes_manager.persistent_attributes = session_attr
    handler_input.attributes_manager.save_persistent_attributes()

    speech_text = "Ok. See you next time!!"

    handler_input.response_builder.speak(speech_text)
    return handler_input.response_builder.response





@sb.request_handler(can_handle_func=lambda input:
                    is_intent_name("AMAZON.FallbackIntent")(input) or
                    is_intent_name("AMAZON.YesIntent")(input) or
                    is_intent_name("AMAZON.NoIntent")(input))
def fallback_handler(handler_input):
    """AMAZON.FallbackIntent is only available in en-US locale.
    This handler will not be triggered except in that locale,
    so it is safe to deploy on any locale.
    """
    # type: (HandlerInput) -> Response
    session_attr = handler_input.attributes_manager.session_attributes

    if ("conversation_state" in session_attr and
            session_attr["conversation_state"]=="STARTED"):
        speech_text = (
            "Sorry {}  can't help you with that.  "
            "Try asking for help from Kal".format(SKILL_NAME))
        reprompt = "Please ask for help from budget buddy"
    else:
        speech_text = (
            "The {} skill can't help you with that.  "
            "It will come up with a number between 0 and 100 and "
            "you try to guess it by saying a number in that range. "
            "Would you like some help from budget buddy?".format(SKILL_NAME))
        reprompt = "Say yes to start communicating with our assistant Kal"

    handler_input.response_builder.speak(speech_text).ask(reprompt)
    return handler_input.response_builder.response


@sb.request_handler(can_handle_func=lambda input: True)
def unhandled_intent_handler(handler_input):
    """Handler for all other unhandled requests."""
    # type: (HandlerInput) -> Response
    speech = "Say yes to continue or no to end the conversation!!"
    handler_input.response_builder.speak(speech).ask(speech)
    return handler_input.response_builder.response


@sb.exception_handler(can_handle_func=lambda i, e: True)
def all_exception_handler(handler_input, exception):
    """Catch all exception handler, log exception and
    respond with custom message.
    """
    # type: (HandlerInput, Exception) -> Response
    logger.error(exception, exc_info=True)
    speech = "Sorry, I can't understand that. Please say again!!"
    handler_input.response_builder.speak(speech).ask(speech)
    return handler_input.response_builder.response


@sb.global_response_interceptor()
def log_response(handler_input, response):
    """Response logger."""
    # type: (HandlerInput, Response) -> None
    logger.info("Response: {}".format(response))


lambda_handler = sb.lambda_handler()
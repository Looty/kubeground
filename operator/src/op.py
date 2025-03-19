"""A Kubernetes controller based on kopf that responsible serving Quests/Checkers for kubeground"""

import os
import logging
import base64
import datetime
import kopf
import kubernetes
import yaml

# Quests
@kopf.on.create('quest.looty.com', 'v1', 'quest')
def on_create(body, patch, **kwargs):
    """On create callback"""

    logging.info("was created!")

@kopf.on.update('quest.looty.com', 'v1', 'quest')
def on_update(body, meta, spec, status, old, new, diff, **kwargs):
    """On update callback"""

    pass

@kopf.on.delete('quest.looty.com', 'v1', 'quest')
def on_delete(body, meta, spec, status, **kwargs):
    """On delete callback"""

    pass

# Checkers
@kopf.on.create('checker.looty.com', 'v1', 'checker')
def on_create(body, patch, **kwargs):
    """On create callback"""

    logging.info("was created!")

@kopf.on.update('checker.looty.com', 'v1', 'checker')
def on_update(body, meta, spec, status, old, new, diff, **kwargs):
    """On update callback"""

    pass

@kopf.on.delete('checker.looty.com', 'v1', 'checker')
def on_delete(body, meta, spec, status, **kwargs):
    """On delete callback"""

    pass

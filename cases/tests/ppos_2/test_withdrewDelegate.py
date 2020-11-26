# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest
from tests.lib.config import EconomicConfig


@pytest.mark.P0
@pytest.mark.compatibility
def test_ROE_001_007_015(client_new_node):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address1)
    assert_code(result, 0)


@pytest.mark.P1
def test_ROE_002(client_new_node):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    cfg = {"gas": 1}
    status = 0
    try:
        client_new_node.delegate.withdrew_delegate(staking_blocknum, address1, transaction_cfg=cfg)
    except BaseException:
        status = 1
    assert status == 1


@pytest.mark.P3
def test_ROE_003(client_new_node):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    illegal_nodeID = "7ee3276fd6b9c7864eb896310b5393324b6db785a2528c00cc28ca8c" \
                     "3f86fc229a86f138b1f1c8e3a942204c03faeb40e3b22ab11b8983c35dc025de42865990"
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address1, node_id=illegal_nodeID)
    log.info(result)
    assert_code(result, 301109)


@pytest.mark.P1
def test_ROE_004(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address1,
                                                        amount=client_new_node.economic.delegate_limit + 1)
    assert_code(result, 301113)


@pytest.mark.P1
def test_ROE_005_018(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    # Return a pledge
    client_new_node.staking.withdrew_staking(address)

    # The next two cycle
    client_new_node.economic.wait_settlement(client_new_node.node, 2)
    amount1 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount1))

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address1)
    assert_code(result, 0)
    amount2 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount2))
    delegate_limit = client_new_node.economic.delegate_limit
    assert delegate_limit - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P1
def test_ROE_006_008(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    value = client_new_node.economic.delegate_limit * 3
    result = client_new_node.delegate.delegate(0, address1, amount=value)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    value = client_new_node.economic.delegate_limit * 2
    amount1 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount1))
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address1, amount=value)
    assert_code(result, 0)
    amount2 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount2))


@pytest.mark.P1
def test_ROE_010(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    client_new_node.economic.env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    lockup_amount = client_new_node.node.web3.toWei(1000, "ether")
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert_code(result, 0)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # create staking
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)

    delegate_amount = client_new_node.node.web3.toWei(500, "ether")
    # Lock account authorization
    result = client_new_node.delegate.delegate(1, address, amount=delegate_amount)
    log.info(result)
    # Own capital account entrustment
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    undelegate_amount = client_new_node.node.web3.toWei(300, "ether")
    log.info("The amount of redemption is greater than the entrustment of the free account")
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    assert undelegate_amount - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P1
def test_ROE_011(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    client_new_node.economic.env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    lockup_amount = client_new_node.node.web3.toWei(1000, "ether")
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert_code(result, 0)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # create staking
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)

    delegate_amount = client_new_node.node.web3.toWei(500, "ether")
    # Lock account authorization
    result = client_new_node.delegate.delegate(1, address, amount=delegate_amount)
    log.info(result)
    # Own capital account entrustment
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    undelegate_amount = client_new_node.node.web3.toWei(700, "ether")
    log.info("The amount of redemption is greater than the entrustment of the free account")
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    assert delegate_amount - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")
    # The next cycle
    client_new_node.economic.wait_settlement(client_new_node.node)
    locked_delegate = delegate_amount - (undelegate_amount - delegate_amount)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # The remaining entrusted amount
    assert msg["Ret"]["Pledge"] == locked_delegate
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    assert amount3 - amount2 == lockup_amount - msg["Ret"]["debt"]


@pytest.mark.P1
def test_ROE_012(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    # create staking
    clinet = client_new_node
    economic = clinet.economic
    node = clinet.node
    address_staking, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    address, _ = economic.account.generate_account(node.web3, economic.delegate_limit * 100)
    result = clinet.staking.create_staking(0, address_staking, address_staking)
    assert_code(result, 0)
    delegate_amount = economic.delegate_limit * 50
    result = clinet.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    amount1 = node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    msg = clinet.ppos.getCandidateInfo(node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    # After redemptive balance is less than the threshold that entrusts gold, redeem completely
    undelegate_amount = node.web3.toWei(49.9, 'ether')
    result = clinet.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)
    time.sleep(2)
    amount2 = node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    result = clinet.ppos.getDelegateInfo(staking_blocknum, address, node.node_id)
    assert_code(result, 301205)
    assert amount1 + delegate_amount - amount2 < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P1
def test_ROE_014(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    lockup_amount = economic.delegate_limit * 1000
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # create staking
    address_staking, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    result = client.staking.create_staking(0, address_staking, address_staking)
    assert_code(result, 0)
    delegate_amount = economic.delegate_limit * 100
    # Lock account authorization
    result = client.delegate.delegate(1, address, amount=delegate_amount)
    assert_code(result, 0)
    # Own capital account entrustment
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)

    msg = client.ppos.getCandidateInfo(node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    undelegate_amount = client_new_node.node.web3.toWei(199.1, "ether")
    amount1 = node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))

    result = client.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)

    result = client.ppos.getDelegateInfo(staking_blocknum, address, node.node_id)
    assert_code(result, 301205)

    amount2 = client.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    assert amount1 + delegate_amount - amount2 < node.web3.toWei(1, "ether")
    #
    # economic.wait_settlement(client_new_node.node)
    # amount3 = client_new_node.node.eth.getBalance(address)
    # log.info("The wallet balance:{}".format(amount3))
    # assert amount3 - amount2 == delegate_amount


@pytest.mark.P1
def test_ROE_017(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    client_new_node.economic.env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    lockup_amount = client_new_node.node.web3.toWei(500, "ether")
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert_code(result, 0)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # create staking
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)

    delegate_amount = client_new_node.node.web3.toWei(500, "ether")
    # Lock account authorization
    result = client_new_node.delegate.delegate(1, address, amount=delegate_amount)
    log.info(result)
    # Own capital account entrustment
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    # Redemptive amount is equal to free account + the entrustment gold of lock storehouse

    undelegate_amount = client_new_node.node.web3.toWei(1000, "ether")
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    assert delegate_amount - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")

    client_new_node.economic.wait_settlement(client_new_node.node)

    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    assert amount3 - amount2 == delegate_amount


@pytest.mark.P1
def test_ROE_019_021(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    value = client_new_node.economic.delegate_limit * 3
    result = client_new_node.delegate.delegate(0, address1, amount=value)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    value = client_new_node.economic.delegate_limit * 2
    amount1 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount1))

    client_new_node.economic.wait_settlement(client_new_node.node)
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address1, amount=value)
    assert_code(result, 0)
    amount2 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount2))
    assert value - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P0
def test_ROE_020(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)

    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    amount1 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount1))

    client_new_node.economic.wait_settlement(client_new_node.node)
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address1)
    assert_code(result, 0)
    amount2 = client_new_node.node.eth.getBalance(address1)
    log.info("The wallet balance:{}".format(amount2))
    delegate_limit = client_new_node.economic.delegate_limit
    assert delegate_limit - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P1
def test_ROE_024(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    client_new_node.economic.env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    lockup_amount = client_new_node.node.web3.toWei(1000, "ether")
    plan = [{'Epoch': 2, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert_code(result, 0)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # create staking
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)

    delegate_amount = client_new_node.node.web3.toWei(500, "ether")
    # Lock account authorization
    result = client_new_node.delegate.delegate(1, address, amount=delegate_amount)
    log.info(result)
    # Own capital account entrustment
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)
    # The next cycle
    client_new_node.economic.wait_settlement(client_new_node.node)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    undelegate_amount = client_new_node.node.web3.toWei(700, "ether")
    log.info("The amount of redemption is greater than the entrustment of the free account")
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    assert delegate_amount - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")

    # The next cycle
    client_new_node.economic.wait_settlement(client_new_node.node)
    locked_delegate = delegate_amount - (undelegate_amount - delegate_amount)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    assert msg["Ret"]["Pledge"] == locked_delegate
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    assert amount3 - amount2 == lockup_amount - msg["Ret"]["debt"]


@pytest.mark.P1
def test_ROE_028(client_new_node):
    """
    :param client_new_node_bgj:
    :return:
    """
    # create staking
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)
    delegate_amount = client_new_node.node.web3.toWei(500, "ether")

    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    # create delegate
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    # The next cycle
    client_new_node.economic.wait_settlement(client_new_node.node)
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=delegate_amount)
    assert_code(result, 0)

    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    assert delegate_amount - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P1
def test_ROE_030(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    client_new_node.economic.env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    lockup_amount = client_new_node.node.web3.toWei(500, "ether")
    plan = [{'Epoch': 2, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert_code(result, 0)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # create staking
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)

    delegate_amount = client_new_node.node.web3.toWei(500, "ether")
    # Lock account authorization
    result = client_new_node.delegate.delegate(1, address, amount=delegate_amount)
    log.info(result)
    # Own capital account entrustment
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)
    # The next cycle
    client_new_node.economic.wait_settlement(client_new_node.node)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    # Redemptive amount is equal to free account + the entrustment gold of lock storehouse

    undelegate_amount = client_new_node.node.web3.toWei(1000, "ether")
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    assert delegate_amount - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")

    client_new_node.economic.wait_settlement(client_new_node.node)

    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    assert amount3 - amount2 == delegate_amount


@pytest.mark.P1
def test_ROE_042(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    client_new_node.economic.env.deploy_all()

    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    lockup_amount = client_new_node.node.web3.toWei(1000, "ether")
    plan = [{'Epoch': 2, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert_code(result, 0)
    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    # create staking
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)

    delegate_amount = client_new_node.node.web3.toWei(200, "ether")
    # Lock account authorization
    result = client_new_node.delegate.delegate(1, address, amount=delegate_amount)
    log.info(result)
    # Own capital account entrustment
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)
    # The next cycle
    client_new_node.economic.wait_settlement(client_new_node.node)

    delegate_amount_2 = client_new_node.node.web3.toWei(300, "ether")
    result = client_new_node.delegate.delegate(1, address, amount=delegate_amount_2)
    log.info(result)
    # Own capital account entrustment
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount_2)
    log.info(result)

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)

    undelegate_amount = client_new_node.node.web3.toWei(700, "ether")
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    log.info("Redemption is more than the hesitation period of the own amount + lock amount")
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=undelegate_amount)
    assert_code(result, 0)

    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address, client_new_node.node.node_id)
    log.info(msg)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    """赎回金额700，先赎回自由金额300，锁仓金额300，再赎回自由金额的100；目前锁定期的自由资金立马退，所以退400"""
    account_dill = undelegate_amount - delegate_amount_2 * 2
    account_sum = account_dill + delegate_amount_2
    assert amount2 - amount1 > account_sum - client_new_node.node.web3.toWei(1, "ether")

    client_new_node.economic.wait_settlement(client_new_node.node)

    msg = client_new_node.ppos.getRestrictingInfo(address)
    log.info(msg)
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    assert amount3 - amount2 == lockup_amount - msg["Ret"]["debt"]


@pytest.mark.P1
def test_ROE_055(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address_staking, address_staking)
    delegate_amount = client_new_node.node.web3.toWei(500, "ether")

    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    # create delegate
    result = client_new_node.delegate.delegate(0, address, amount=delegate_amount)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, amount=delegate_amount + 1)
    assert_code(result, 301113)


@pytest.mark.P1
def test_ROE_056_057(client_new_node, client_consensus):
    """

    :param client_new_node_obj:
    :param client_consensus_obj:
    :return:
    """
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    value = client_new_node.economic.create_staking_limit * 2
    result = client_new_node.staking.create_staking(0, address_staking, address_staking,
                                                    amount=value)
    assert_code(result, 0)

    # create delegate
    result = client_new_node.delegate.delegate(0, address)
    log.info(result)

    log.info("Close one node")
    client_new_node.node.stop()
    node = client_consensus.node

    msg = client_consensus.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    log.info("The next two periods")
    client_consensus.economic.wait_settlement(node, 2)

    log.info("Restart the node")
    client_new_node.node.start()
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address)
    log.info(result)

    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    delegate_limit = client_new_node.economic.delegate_limit
    assert delegate_limit - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P3
def test_ROE_058(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address_staking, address_staking)
    assert_code(result, 0)

    result = client_new_node.delegate.delegate(0, address)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    amount = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount))
    cfg = {"gasPrice": amount}
    status = 0
    try:
        result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, transaction_cfg=cfg)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@pytest.mark.P3
def test_ROE_059(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address_staking, address_staking)
    assert_code(result, 0)

    result = client_new_node.delegate.delegate(0, address)
    log.info(result)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    cfg = {"gas": 10}
    status = 0
    try:
        result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address, transaction_cfg=cfg)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@pytest.mark.P1
def test_ROE_060(client_new_node):
    """

    :param client_new_node_obj:
    :return:
    """
    address_staking, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           10 ** 18 * 10000000)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    # The next cycle
    client_new_node.economic.wait_settlement(client_new_node.node)

    result = client_new_node.staking.create_staking(0, address_staking, address_staking)
    assert_code(result, 0)

    result = client_new_node.staking.withdrew_staking(address_staking)
    assert_code(result, 0)
    # The next two cycle
    client_new_node.economic.wait_settlement(client_new_node.node, 2)
    # Pledge again after quitting pledge
    result = client_new_node.staking.create_staking(0, address_staking, address_staking)
    assert_code(result, 0)
    result = client_new_node.delegate.delegate(0, address)
    assert_code(result, 0)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address)
    log.info(result)
    assert_code(result, 0)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    delegate_limit = client_new_node.economic.delegate_limit
    assert delegate_limit - (amount2 - amount1) < client_new_node.node.web3.toWei(1, "ether")

-- phpMyAdmin SQL Dump
-- version 5.0.2
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Jul 23, 2020 at 04:39 AM
-- Server version: 10.4.11-MariaDB
-- PHP Version: 7.4.6

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `mam_core`
--

-- --------------------------------------------------------

--
-- Table structure for table `ffs_alloc_instrument`
--

CREATE TABLE `ffs_alloc_instrument` (
  `alloc_instrument_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `periode_key` int(11) NOT NULL,
  `instrument_key` int(11) NOT NULL,
  `instrument_value` decimal(8,4) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ffs_alloc_sector`
--

CREATE TABLE `ffs_alloc_sector` (
  `alloc_sector_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `periode_key` int(11) NOT NULL,
  `sector_key` int(11) NOT NULL,
  `sector_value` decimal(8,4) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ffs_alloc_security`
--

CREATE TABLE `ffs_alloc_security` (
  `alloc_security_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `periode_key` int(11) NOT NULL,
  `sec_key` int(11) NOT NULL,
  `security_value` decimal(8,4) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ffs_nav_performance`
--

CREATE TABLE `ffs_nav_performance` (
  `nav_perform_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `periode_key` int(11) NOT NULL,
  `nav_date` datetime(6) NOT NULL,
  `nav_d0` decimal(9,4) NOT NULL,
  `nav_d1` decimal(9,4) NOT NULL,
  `nav_m0` decimal(9,4) NOT NULL,
  `nav_m1` decimal(9,4) NOT NULL,
  `nav_m3` decimal(9,4) NOT NULL,
  `nav_m6` decimal(9,4) NOT NULL,
  `nav_ytd` decimal(9,4) NOT NULL,
  `nav_y1` decimal(9,4) NOT NULL,
  `nav_y3` decimal(9,4) NOT NULL,
  `nav_y5` decimal(9,4) NOT NULL,
  `perform_d1` decimal(9,4) NOT NULL,
  `perform_mtd` decimal(9,4) NOT NULL,
  `perform_m1` decimal(9,4) NOT NULL,
  `perform_m3` decimal(9,4) NOT NULL,
  `perform_m6` decimal(9,4) NOT NULL,
  `perfrom_ytd` decimal(9,4) NOT NULL,
  `perform_y1` decimal(9,4) NOT NULL,
  `perform_y3` decimal(9,4) NOT NULL,
  `perform_y5` decimal(9,4) NOT NULL,
  `perform_cagr` decimal(9,4) NOT NULL,
  `perform_all` decimal(9,4) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ffs_periode`
--

CREATE TABLE `ffs_periode` (
  `periode_key` int(11) NOT NULL,
  `periode_date` date NOT NULL,
  `periode_name` varchar(50) NOT NULL,
  `date_opened` date DEFAULT NULL,
  `date_closed` date DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ffs_publish`
--

CREATE TABLE `ffs_publish` (
  `ffs_key` int(11) NOT NULL,
  `periode_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `ffs_link` varchar(255) DEFAULT NULL,
  `date_periode` date NOT NULL,
  `date_published` date DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent`
--

CREATE TABLE `ms_agent` (
  `agent_key` int(11) NOT NULL,
  `agent_id` int(11) NOT NULL,
  `agent_code` varchar(20) NOT NULL,
  `agent_name` varchar(100) NOT NULL,
  `agent_short_name` varchar(50) DEFAULT NULL,
  `agent_category` int(11) NOT NULL,
  `agent_channel` int(11) NOT NULL,
  `reference_code` varchar(50) DEFAULT NULL,
  `remarks` varchar(100) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_agent`
--

INSERT INTO `ms_agent` (`agent_key`, `agent_id`, `agent_code`, `agent_name`, `agent_short_name`, `agent_category`, `agent_channel`, `reference_code`, `remarks`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 173, '173', 'SALES HO NO_SALES(POOL)', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 132, '132', 'SALES HO MAM-GROUP', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 137, '137', 'SALES HO MAM-NON_GROUP', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 303, '303', 'SALES MNC DUIT', NULL, 0, 2, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 133, '133', 'SALES BANDUNG(POOL)', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 376, '376', 'SALES MOHAMAD RAZAB', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(7, 382, '382', 'SALES NITA RAHAYU UTAMI', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(8, 355, '355', 'SALES DIMAS', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(9, 265, '265', 'SALES REZHA', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(10, 266, '266', 'SALES BURHAN', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(11, 271, '271', 'SALES RIMA', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(12, 351, '351', 'SALES WISYNU FARANDIKHA', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(13, 353, '353', 'SALES PEKANBARU(POOL)', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(14, 390, '390', 'SALES CHINTYA JONETA', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(15, 228, '228', 'SALES BIC-SURABAYA', NULL, 1, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(16, 327, '327', 'SALES MEDAN(POOL)', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(17, 340, '340', 'SALES LEO', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(18, 335, '335', 'SALES BUDI RIMPAN', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(19, 350, '350', 'SALES HENGKI', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(20, 356, '356', 'SALES PRISSILLIA FRANSISKA', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(21, 339, '339', 'SALES DINI', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(22, 363, '363', 'SALES REFERAL (ANNAFRID)', NULL, 0, 1, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(23, 352, '352', 'SALES AGUNG H EKA PUTRA', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(24, 374, '374', 'SALES BAGUS RADITYO', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(25, 307, '307', 'SALES ANISA DWI OLIVIA', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(26, 378, '378', 'SALES YOSPIN', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(27, 163, '163', 'SALES ENITA', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(28, 306, '306', 'SALES KURNIA', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(29, 309, '309', 'SALES DHEKHA', NULL, 0, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(30, 230, 'BRI', 'SALES APERD BRI', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(31, 357, 'BNI', 'SALES APERD BNI', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(32, 229, 'BAREKSA', 'SALES APERD-BUANA CAPITAL(BAREKSA)', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(33, 231, 'IPOT', 'SALES APERD-IPOT', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(34, 232, 'PHILLIP', 'SALES APERD-PHILLIP SECURITIES', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(35, 279, 'MNC69', 'SALES APERD-MSEC', NULL, 2, 2, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(36, 358, 'MIRAE', 'SALES APERD-MIRAE', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(37, 366, 'INVESTAMART', 'SALES APERD-INVESTAMART(ACTIVE)', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(38, 381, 'TAKJUB', 'SALES APERD TAKJUB TEKNOLOGI INDONESIA', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(39, 385, 'MANDIRI', 'SALES APERD MANDIRI SEKURITAS', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(40, 389, 'BNISEK', 'SALES APERD BNI SEKURITAS', NULL, 2, 0, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent_agreement`
--

CREATE TABLE `ms_agent_agreement` (
  `agreement_key` int(11) NOT NULL,
  `agreement_no` varchar(20) DEFAULT NULL,
  `agreement_date` date NOT NULL,
  `agreement_subject` varchar(100) NOT NULL,
  `agreement_content` varchar(1000) DEFAULT NULL,
  `signed_date` date DEFAULT NULL,
  `sign_city` varchar(50) DEFAULT NULL,
  `branch_key` int(11) DEFAULT NULL,
  `agreement_status` int(11) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent_bank_account`
--

CREATE TABLE `ms_agent_bank_account` (
  `agent_bankacc_key` int(11) NOT NULL,
  `agent_key` int(11) NOT NULL,
  `bank_account_seqno` int(11) NOT NULL,
  `bank_account_key` int(11) NOT NULL,
  `bank_account_purpose` varchar(50) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent_branch`
--

CREATE TABLE `ms_agent_branch` (
  `agent_branch_key` int(11) NOT NULL,
  `agent_key` int(11) NOT NULL,
  `branch_key` int(11) NOT NULL,
  `eff_date` date DEFAULT NULL,
  `remarks` longtext DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_agent_branch`
--

INSERT INTO `ms_agent_branch` (`agent_branch_key`, `agent_key`, `branch_key`, `eff_date`, `remarks`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 1, 2, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 2, 2, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 3, 2, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 4, 2, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 35, 19, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent_customer`
--

CREATE TABLE `ms_agent_customer` (
  `agent_customer_key` int(11) NOT NULL,
  `agent_key` int(11) NOT NULL,
  `customer_key` int(11) NOT NULL,
  `eff_date` date DEFAULT NULL,
  `ref_code_used` varchar(20) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_agent_customer`
--

INSERT INTO `ms_agent_customer` (`agent_customer_key`, `agent_key`, `customer_key`, `eff_date`, `ref_code_used`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 2, 1, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 32, 2, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 30, 3, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 33, 4, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 34, 5, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 35, 6, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(7, 31, 7, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(8, 36, 8, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(9, 37, 9, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(10, 38, 10, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(11, 39, 11, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(12, 40, 12, '2020-01-01', NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent_detail`
--

CREATE TABLE `ms_agent_detail` (
  `agent_key` int(11) NOT NULL,
  `mobile_no` int(11) NOT NULL,
  `mobile_no_alt` int(11) DEFAULT NULL,
  `email_address` varchar(50) NOT NULL,
  `birth_date` datetime(6) DEFAULT NULL,
  `birth_place` varchar(50) DEFAULT NULL,
  `gender` int(11) NOT NULL,
  `nationality` int(11) NOT NULL,
  `country_key` int(11) DEFAULT NULL,
  `IDType` int(11) NOT NULL,
  `id_no` varchar(50) DEFAULT NULL,
  `occupation` varchar(50) DEFAULT NULL,
  `company_name` varchar(50) DEFAULT NULL,
  `position` varchar(50) DEFAULT NULL,
  `position_level` int(11) NOT NULL,
  `join_date` datetime(6) NOT NULL,
  `flag_resign` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent_license`
--

CREATE TABLE `ms_agent_license` (
  `aglic_key` int(11) NOT NULL,
  `agent_key` int(11) NOT NULL,
  `license_name` varchar(50) DEFAULT NULL,
  `license_no` varchar(50) DEFAULT NULL,
  `license_issuer` varchar(50) DEFAULT NULL,
  `license_exp_date` date DEFAULT NULL,
  `alert_before_expired` tinyint(1) NOT NULL,
  `license_notes` varchar(150) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_agent_product`
--

CREATE TABLE `ms_agent_product` (
  `agent_product_key` int(11) NOT NULL,
  `branch_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `eff_date` date DEFAULT NULL,
  `product_name_sa` varchar(50) DEFAULT NULL,
  `holding_period_days` int(11) NOT NULL,
  `mgt_fee_share_sa` decimal(5,2) NOT NULL,
  `sub_fee_share_sa` decimal(5,2) NOT NULL,
  `red_fee_share_sa` decimal(5,2) NOT NULL,
  `swtot_fee_share_sa` decimal(5,2) NOT NULL,
  `swtin_fee_share_sa` decimal(5,2) NOT NULL,
  `ojk_fee_share_sa` decimal(5,2) NOT NULL,
  `other_fee_share_sa` decimal(5,2) NOT NULL,
  `flag_enabled` tinyint(1) NOT NULL,
  `flag_subscription` tinyint(1) NOT NULL,
  `flag_redemption` tinyint(1) NOT NULL,
  `flag_switch_out` tinyint(1) NOT NULL,
  `flag_switch_in` tinyint(1) NOT NULL,
  `max_sub_fee` decimal(12,4) NOT NULL,
  `max_red_fee` decimal(12,4) NOT NULL,
  `max_swi_fee` decimal(12,4) NOT NULL,
  `min_sub_amount` decimal(12,4) NOT NULL,
  `min_red_amount` decimal(12,4) NOT NULL,
  `min_red_unit` decimal(12,4) NOT NULL,
  `min_unit_after_red` decimal(12,4) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_bank`
--

CREATE TABLE `ms_bank` (
  `bank_key` int(11) NOT NULL,
  `bank_code` varchar(20) NOT NULL,
  `bank_name` varchar(150) NOT NULL,
  `bank_fullname` varchar(150) DEFAULT NULL,
  `bi_member_code` varchar(20) DEFAULT NULL,
  `swift_code` varchar(20) DEFAULT NULL,
  `flag_local` tinyint(1) NOT NULL,
  `flag_government` tinyint(1) NOT NULL,
  `bank_logo` varchar(255) DEFAULT NULL,
  `bank_ibank_url` varchar(255) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_bank_account`
--

CREATE TABLE `ms_bank_account` (
  `bank_account_key` int(11) NOT NULL,
  `bank_key` int(11) NOT NULL,
  `account_no` varchar(20) NOT NULL,
  `account_holder_name` varchar(50) NOT NULL,
  `branch_name` varchar(50) DEFAULT NULL,
  `currency_key` int(11) NOT NULL,
  `bank_account_type` int(11) NOT NULL,
  `rec_domain` int(11) NOT NULL,
  `swift_code` varchar(20) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_branch`
--

CREATE TABLE `ms_branch` (
  `branch_key` int(11) NOT NULL,
  `participant_key` int(11) DEFAULT NULL,
  `branch_code` varchar(30) NOT NULL,
  `branch_name` varchar(100) NOT NULL,
  `branch_category` varchar(10) NOT NULL,
  `city_key` int(11) DEFAULT NULL,
  `branch_address` varchar(150) DEFAULT NULL,
  `branch_established` datetime(6) DEFAULT NULL,
  `branch_pic_name` varchar(50) DEFAULT NULL,
  `branch_pic_email` varchar(50) DEFAULT NULL,
  `branch_pic_phoneno` varchar(20) DEFAULT NULL,
  `branch_cost_center` varchar(10) DEFAULT NULL,
  `branch_profit_center` varchar(10) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_branch`
--

INSERT INTO `ms_branch` (`branch_key`, `participant_key`, `branch_code`, `branch_name`, `branch_category`, `city_key`, `branch_address`, `branch_established`, `branch_pic_name`, `branch_pic_email`, `branch_pic_phoneno`, `branch_cost_center`, `branch_profit_center`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, NULL, 'HOPOOL', 'JAKARTA POOL (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '9', NULL, NULL),
(2, NULL, 'HOMAM', 'MNC ASSET HO (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '321', NULL, NULL),
(3, NULL, 'KBS', 'JAKARTA-KBS (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '331', NULL, NULL),
(4, NULL, 'SYARI', 'JAKARTA-SYARIAH (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '332', NULL, NULL),
(5, NULL, 'KBS2', 'JAKARTA-KEBON SIRIH-2 (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '25', NULL, NULL),
(6, NULL, 'INSTI', 'JAKARTA-INSTITUSI (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '273', NULL, NULL),
(7, NULL, 'KBS1', 'JAKARTA-KEBON SIRIH-1 (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '285', NULL, NULL),
(8, NULL, 'ETF', 'JAKARTA-ETF (BRANCH)', 'IM_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '287', NULL, NULL),
(9, NULL, 'BGR', 'BOGOR (BRANCH)', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '286', NULL, NULL),
(10, NULL, 'BDG', 'BANDUNG (BRANCH)', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '6', NULL, NULL),
(11, NULL, 'SBY', 'SURABAYA (BRANCH)', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '15', NULL, NULL),
(12, NULL, 'BIC', 'SURABAYA-BIC (BRANCH)', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '320', NULL, NULL),
(13, NULL, 'MDN', 'MEDAN (BRANCH)', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '324', NULL, NULL),
(14, NULL, 'PKU', 'PEKANBARU (BRANCH)', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '315', NULL, NULL),
(15, NULL, 'BAREKSA', 'APERD BAREKSA (BUANA CAPITAL)', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '19', NULL, NULL),
(16, NULL, 'BRI', 'APERD BRI', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '20', NULL, NULL),
(17, NULL, 'PHILLIP', 'APERD PHILLIP SEKURITAS', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '21', NULL, NULL),
(18, NULL, 'IPOD', 'APERD IPOT FUND', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '22', NULL, NULL),
(19, NULL, 'MNC69', 'APERD MNC SEKURITAS', 'SA_BRANCH', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '280', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_city`
--

CREATE TABLE `ms_city` (
  `city_key` int(11) NOT NULL,
  `parent_key` int(11) DEFAULT NULL,
  `city_code` varchar(12) NOT NULL,
  `city_name` varchar(50) NOT NULL,
  `city_level` int(11) NOT NULL,
  `postal_code` varchar(6) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_country`
--

CREATE TABLE `ms_country` (
  `country_key` int(11) NOT NULL,
  `cou_code` varchar(5) NOT NULL,
  `cou_name` varchar(50) NOT NULL,
  `short_name` varchar(30) DEFAULT NULL,
  `flag_base` tinyint(1) NOT NULL,
  `currency_key` int(11) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_country`
--

INSERT INTO `ms_country` (`country_key`, `cou_code`, `cou_name`, `short_name`, `flag_base`, `currency_key`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'AD', 'Andorra', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'AE', 'United Arab Emirates', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 'AF', 'Afghanistan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 'AG', 'Antigua And Barbuda', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 'AI', 'Anguilla', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 'AL', 'Albania', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(7, 'AM', 'Armenia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(8, 'AN', 'Netherlands Antilles', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(9, 'AO', 'Angola', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(10, 'AQ', 'Antarctica', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(11, 'AR', 'Argentina', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(12, 'AS', 'American Samoa', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(13, 'AT', 'Austria', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(14, 'AU', 'Australia', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(15, 'AW', 'Aruba', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(16, 'AZ', 'Azerbaijan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(17, 'BA', 'Bosnia And Herzegovina', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(18, 'BB', 'Barbados', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(19, 'BD', 'Bangladesh', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(20, 'BE', 'Belgium', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(21, 'BF', 'Burkina Faso', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(22, 'BG', 'Bulgaria', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(23, 'BH', 'Bahrain', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(24, 'BI', 'Burundi', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(25, 'BJ', 'Benin', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(26, 'BM', 'Bermuda', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(27, 'BN', 'Brunei Darussalam', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(28, 'BO', 'Bolivia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(29, 'BR', 'Brazil', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(30, 'BS', 'Bahamas', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(31, 'BT', 'Bhutan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(32, 'BV', 'Bouvet Island', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(33, 'BW', 'Botswana', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(34, 'BY', 'Belarus', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(35, 'BZ', 'Belize', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(36, 'CA', 'Canada', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(37, 'CC', 'Cocos(keeling) Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(38, 'CD', 'Congo, The Democratic Republic Of The', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(39, 'CF', 'Central African Republic', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(40, 'CG', 'Congo', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(41, 'CH', 'Switzerland', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(42, 'CI', 'Cote D\'ivoire', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(43, 'CK', 'Cook Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(44, 'CL', 'Chile', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(45, 'CM', 'Cameroon', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(46, 'CN', 'China', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(47, 'CO', 'Colombia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(48, 'CR', 'Costa Rica', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(49, 'CU', 'Cuba', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(50, 'CV', 'Cape Verde', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(51, 'CX', 'Christmas Island', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(52, 'CY', 'Cyprus', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(53, 'CZ', 'Czech Republic', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(54, 'DE', 'Germany', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(55, 'DJ', 'Djibouti', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(56, 'DK', 'Denmark', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(57, 'DM', 'Dominica', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(58, 'DO', 'Dominican Republic', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(59, 'DZ', 'Algeria', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(60, 'EC', 'Ecuador', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(61, 'EE', 'Estonia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(62, 'EG', 'Egypt', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(63, 'EH', 'Western Sahara', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(64, 'ER', 'Eritrea', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(65, 'ES', 'Spain', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(66, 'ET', 'Ethiopia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(67, 'FI', 'Finland', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(68, 'FJ', 'Fiji', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(69, 'FK', 'Falkland Islands (malvinas)', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(70, 'FM', 'Micronesia, Federated States Of', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(71, 'FO', 'Faroe Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(72, 'FR', 'France', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(73, 'GA', 'Gabon', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(74, 'GB', 'United Kingdom', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(75, 'GD', 'Grenada', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(76, 'GE', 'Georgia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(77, 'GF', 'French Guiana', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(78, 'GH', 'Ghana', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(79, 'GI', 'Gibraltar', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(80, 'GL', 'Greenland', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(81, 'GM', 'Gambia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(82, 'GN', 'Guinea', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(83, 'GP', 'Guadeloupe', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(84, 'GQ', 'Equatorial Guinea', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(85, 'GR', 'Greece', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(86, 'GS', 'South Georgia And The South Sandwich Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(87, 'GT', 'Guatemala', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(88, 'GU', 'Guam', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(89, 'GW', 'Guinea-bissau', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(90, 'GY', 'Guyana', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(91, 'HK', 'Hong Kong', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(92, 'HM', 'Heard Island And Mcdonald Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(93, 'HN', 'Honduras', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(94, 'HR', 'Croatia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(95, 'HT', 'Haiti', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(96, 'HU', 'Hungary', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(97, 'ID', 'Indonesia', NULL, 1, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(98, 'IE', 'Ireland', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(99, 'IL', 'Israel', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(100, 'IN', 'India', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(101, 'IO', 'British Indian Ocean Territory', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(102, 'IQ', 'Iraq', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(103, 'IR', 'Iran, Islamic Republic Of', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(104, 'IS', 'Iceland', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(105, 'IT', 'Italy', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(106, 'JM', 'Jamaica', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(107, 'JO', 'Jordan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(108, 'JP', 'Japan', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(109, 'KE', 'Kenya', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(110, 'KG', 'Kyrgyzstan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(111, 'KH', 'Cambodia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(112, 'KI', 'Kiribati', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(113, 'KM', 'Comoros', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(114, 'KN', 'Saint Kitts And Nevis', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(115, 'KP', 'Korea, Democratic Peoples Republic Of', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(116, 'KR', 'Korea, Republic Of', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(117, 'KV', 'Kosovo', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(118, 'KW', 'Kuwait', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(119, 'KY', 'Cayman Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(120, 'KZ', 'Kazakstan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(121, 'LA', 'Lao Peoples Democratic Republic', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(122, 'LB', 'Lebanon', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(123, 'LC', 'Saint Lucia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(124, 'LI', 'Liechtenstein', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(125, 'LK', 'Sri Lanka', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(126, 'LR', 'Liberia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(127, 'LS', 'Lesotho', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(128, 'LT', 'Lithuania', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(129, 'LU', 'Luxembourg', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(130, 'LV', 'Latvia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(131, 'LY', 'Libyan Arab Jamahiriya', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(132, 'MA', 'Morocco', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(133, 'MC', 'Monaco', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(134, 'MD', 'Moldova, Republic Of', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(135, 'ME', 'Montenegro', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(136, 'MG', 'Madagascar', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(137, 'MH', 'Marshall Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(138, 'MK', 'Macedonia, The Former Yugoslav Republic Of', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(139, 'ML', 'Mali', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(140, 'MM', 'Myanmar', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(141, 'MN', 'Mongolia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(142, 'MO', 'Macau', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(143, 'MP', 'Northern Mariana Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(144, 'MQ', 'Martinique', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(145, 'MR', 'Mauritania', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(146, 'MS', 'Montserrat', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(147, 'MT', 'Malta', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(148, 'MU', 'Mauritius', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(149, 'MV', 'Maldives', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(150, 'MW', 'Malawi', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(151, 'MX', 'Mexico', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(152, 'MY', 'Malaysia', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(153, 'MZ', 'Mozambique', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(154, 'NA', 'Namibia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(155, 'NC', 'New Caledonia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(156, 'NE', 'Niger', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(157, 'NF', 'Norfolk Island', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(158, 'NG', 'Nigeria', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(159, 'NI', 'Nicaragua', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(160, 'NL', 'Netherlands', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(161, 'NO', 'Norway', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(162, 'NP', 'Nepal', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(163, 'NR', 'Nauru', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(164, 'NU', 'Niue', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(165, 'NZ', 'New Zealand', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(166, 'OM', 'Oman', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(167, 'PA', 'Panama', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(168, 'PE', 'Peru', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(169, 'PF', 'French Polynesia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(170, 'PG', 'Papua New Guinea', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(171, 'PH', 'Philippines', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(172, 'PK', 'Pakistan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(173, 'PL', 'Poland', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(174, 'PM', 'Saint Pierre And Miquelon', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(175, 'PN', 'Pitcairn', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(176, 'PR', 'Puerto Rico', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(177, 'PS', 'Palestinian Territory, Occupied', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(178, 'PT', 'Portugal', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(179, 'PW', 'Palau', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(180, 'PY', 'Paraguay', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(181, 'QA', 'Qatar', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(182, 'RE', 'Reunion', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(183, 'RO', 'Romania', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(184, 'RS', 'Serbia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(185, 'RU', 'Russian Federation', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(186, 'RW', 'Rwanda', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(187, 'SA', 'Saudi Arabia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(188, 'SB', 'Solomon Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(189, 'SC', 'Seychelles', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(190, 'SD', 'Sudan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(191, 'SE', 'Sweden', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(192, 'SG', 'Singapore', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(193, 'SH', 'Saint Helena', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(194, 'SI', 'Slovenia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(195, 'SJ', 'Svalbard And Jan Mayen', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(196, 'SK', 'Slovakia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(197, 'SL', 'Sierra Leone', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(198, 'SM', 'San Marino', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(199, 'SN', 'Senegal', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(200, 'SO', 'Somalia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(201, 'SR', 'Suriname', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(202, 'ST', 'Sao Tome And Principe', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(203, 'SV', 'El Salvador', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(204, 'SY', 'Syrian Arab Republic', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(205, 'SZ', 'Swaziland', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(206, 'TC', 'Turks And Caicos Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(207, 'TD', 'Chad', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(208, 'TF', 'French Southern Territories', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(209, 'TG', 'Togo', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(210, 'TH', 'Thailand', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(211, 'TJ', 'Tajikistan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(212, 'TK', 'Tokelau', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(213, 'TM', 'Turkmenistan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(214, 'TN', 'Tunisia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(215, 'TO', 'Tonga', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(216, 'TP', 'East Timor', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(217, 'TR', 'Turkey', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(218, 'TT', 'Trinidad And Tobago', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(219, 'TV', 'Tuvalu', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(220, 'TW', 'Taiwan, Province Of China', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(221, 'TZ', 'Tanzania, United Republic Of', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(222, 'UA', 'Ukraine', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(223, 'UG', 'Uganda', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(224, 'UM', 'United States Minor Outlying Islands', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(225, 'US', 'United States', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(226, 'UY', 'Uruguay', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(227, 'UZ', 'Uzbekistan', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(228, 'VA', 'Holy See(Vatican City State)', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(229, 'VC', 'Saint Vincent And The Grenadines', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(230, 'VE', 'Venezuela', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(231, 'VG', 'Virgin Islands, British', NULL, 0, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(232, 'VI', 'Virgin Islands, U.s.', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(233, 'VN', 'Viet Nam', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(234, 'VU', 'Vanuatu', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(235, 'WF', 'Wallis And Futuna', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(236, 'WS', 'Samoa', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(237, 'YE', 'Yemen', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(238, 'YT', 'Mayotte', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(239, 'ZA', 'South Africa', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(240, 'ZM', 'Zambia', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(241, 'ZW', 'Zimbabwe', NULL, 0, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_currency`
--

CREATE TABLE `ms_currency` (
  `currency_key` int(11) NOT NULL,
  `code` varchar(3) NOT NULL,
  `symbol` varchar(3) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `flag_base` tinyint(1) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_currency`
--

INSERT INTO `ms_currency` (`currency_key`, `code`, `symbol`, `name`, `flag_base`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'IDR', 'Rp', 'Indonesia Rupiah', 1, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '1', NULL, NULL),
(2, 'USD', '$', 'US Dollar', 1, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_custodian_bank`
--

CREATE TABLE `ms_custodian_bank` (
  `custodian_key` int(11) NOT NULL,
  `custodian_code` varchar(20) NOT NULL,
  `custodian_short_name` varchar(30) NOT NULL,
  `custodian_full_name` varchar(50) DEFAULT NULL,
  `bi_member_code` varchar(20) DEFAULT NULL,
  `swift_code` varchar(20) DEFAULT NULL,
  `flag_local` tinyint(1) NOT NULL,
  `flag_government` tinyint(1) NOT NULL,
  `bank_logo` varchar(255) DEFAULT NULL,
  `custodian_profile` varchar(1000) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_custodian_bank`
--

INSERT INTO `ms_custodian_bank` (`custodian_key`, `custodian_code`, `custodian_short_name`, `custodian_full_name`, `bi_member_code`, `swift_code`, `flag_local`, `flag_government`, `bank_logo`, `custodian_profile`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'BCA', 'Bank BCA', 'Bank BCA', NULL, NULL, 0, 0, NULL, 'PT Bank Central Asia Tbk (BCA) mulai beroperasi sebagai  Bank Umum berdasarkan Nomor Izin Usaha 42855/U.M.II. Aktivitas BCA sebagai Bank Kustodian dimulai sejak memiliki Izin Usaha BK berdasarkan SK Ketua Bapepam-LK Nomor KEP-148/PM/1991 tanggal 13 November 1991 dengan Kode BK BCA01.', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '12', NULL, NULL),
(2, 'BRI', 'Bank BRI', 'Bank BRI', NULL, NULL, 0, 0, NULL, 'PT Bank Rakyat Indonesia (Persero) Tbk (BRI) mulai beroperasi sebagai Bank Umum berdasarkan Nomor Izin Usaha PP No.21 Tahun 1992 tanggal 29 April 1992. Aktivitas BRI sebagai Bank Kustodian dimulai sejak memiliki Izin Usaha BK berdasarkan SK Ketua Bapepam-LK Nomor KEP-91/PM/1996 Tanggal  11 April 1996 dengan Kode BK BRI01.', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '35', NULL, NULL),
(3, 'BNI', 'Bank BNI', 'Bank BNI', NULL, NULL, 0, 0, NULL, 'PT Bank Negara Indonesia (Persero) Tbk (BNI) mulai beroperasi sebagai Bank Umum berdasarkan  Nomor Izin Usaha UU RI No.17/1968 ttg Bank Negara Indonesia 1946. Aktivitas BNI sebagai Bank Kustodian dimulai sejak memiliki Izin Usaha BK berdasarkan SK Ketua Bapepam-LK Nomor KEP-162/PM/1991 Tanggal 09 Desember 1991.', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '81', NULL, NULL),
(4, 'MAYBANK', 'Maybank', 'Bank Maybank', NULL, NULL, 0, 0, NULL, 'PT Bank Maybank Indonesia Tbk (Maybank) mulai beroperasi sebagai Bank Umum berdasarkan  Nomor Izin Usaha 1384.12/U.M.II. Aktivitas Maybank sebagai Bank Kustodian dimulai sejak memiliki Izin Usaha BK berdasarkan SK Ketua Bapepam-LK Nomor KEP-67/PM/1991 Tanggal 30 Juli 1991.', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_customer`
--

CREATE TABLE `ms_customer` (
  `customer_key` int(11) NOT NULL,
  `id_customer` int(11) NOT NULL,
  `unit_holder_idno` varchar(20) NOT NULL,
  `full_name` varchar(150) NOT NULL,
  `sid_no` varchar(20) DEFAULT NULL,
  `investor_type` varchar(10) NOT NULL,
  `customer_category` varchar(15) NOT NULL,
  `participant_key` int(11) DEFAULT NULL,
  `cif_suspend_flag` tinyint(1) NOT NULL,
  `cif_suspend_modified_date` datetime(6) DEFAULT NULL,
  `cif_suspend_reason` varchar(100) DEFAULT NULL,
  `openacc_branch_key` int(11) DEFAULT NULL,
  `openacc_agent_key` int(11) DEFAULT NULL,
  `openacc_date` datetime(6) DEFAULT NULL,
  `closeacc_branch_key` int(11) DEFAULT NULL,
  `closeacc_agent_key` int(11) DEFAULT NULL,
  `closeacc_date` datetime(6) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_customer`
--

INSERT INTO `ms_customer` (`customer_key`, `id_customer`, `unit_holder_idno`, `full_name`, `sid_no`, `investor_type`, `customer_category`, `participant_key`, `cif_suspend_flag`, `cif_suspend_modified_date`, `cif_suspend_reason`, `openacc_branch_key`, `openacc_agent_key`, `openacc_date`, `closeacc_branch_key`, `closeacc_agent_key`, `closeacc_date`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 3752, '201303000074', 'PT MNC ASSET MANAGEMENT', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '3752', NULL, NULL),
(2, 8360, '201506000003', 'PT BUANA CAPITAL', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '8360', NULL, NULL),
(3, 49919, '201707000073', 'PT BANK BRI (PERSERO) TBK', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '49919', NULL, NULL),
(4, 8022, '201408000040', 'PT INDO PREMIER SECURITIES', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '8022', NULL, NULL),
(5, 8275, '201504000017', 'PT PHILLIP SECURITIES INDONESIA', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '8275', NULL, NULL),
(6, 18826, '201604000097', 'PT MNC SEKURITAS', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '18826', NULL, NULL),
(7, 120710, '201804000039', 'PT BANK NEGARA INDONESIA (PERSERO) TBK', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '120710', NULL, NULL),
(8, 120756, '201805000003', 'PT MIRAE ASSET SEKURITAS INDONESIA', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '120756', NULL, NULL),
(9, 130838, '201806000014', 'PT INVESTAMART PRINCIPAL OPTIMA', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '130838', NULL, NULL),
(10, 151364, '201905000044', 'PT TAKJUB TEKNOLOGI INDONESIA', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '151364', NULL, NULL),
(11, 151395, '201906000031', 'PT MANDIRI SEKURITAS', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '151395', NULL, NULL),
(12, 151442, '201907000046', 'PT BNI SEKURITAS', NULL, 'CUST_INDI', 'DIRECT_CUST', NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '151442', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_customer_login`
--

CREATE TABLE `ms_customer_login` (
  `cust_login_key` int(11) NOT NULL,
  `customer_key` int(11) NOT NULL,
  `login_userid` varchar(30) NOT NULL,
  `login_user_password` varchar(30) NOT NULL,
  `login_user_name` varchar(50) NOT NULL,
  `login_pinno` varchar(50) DEFAULT NULL,
  `login_active` tinyint(1) NOT NULL,
  `email_address` varchar(50) DEFAULT NULL,
  `email_date_verified` datetime(6) DEFAULT NULL,
  `mobile_no` varchar(20) DEFAULT NULL,
  `mobileno_date_verified` datetime(6) DEFAULT NULL,
  `login_category` int(11) NOT NULL,
  `login_customer_type` int(11) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_cutomer_detail`
--

CREATE TABLE `ms_cutomer_detail` (
  `customer_key` int(11) NOT NULL,
  `nationality` int(11) NOT NULL,
  `gender` int(11) NOT NULL,
  `id_type` int(11) NOT NULL,
  `id_number` varchar(20) DEFAULT NULL,
  `id_holder_name` varchar(50) DEFAULT NULL,
  `flag_employee` tinyint(1) NOT NULL,
  `flag_group` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_file`
--

CREATE TABLE `ms_file` (
  `file_key` int(11) NOT NULL,
  `ref_fk_key` int(11) NOT NULL,
  `ref_fk_domain` varchar(30) NOT NULL,
  `file_name` varchar(50) NOT NULL,
  `file_ext` varchar(6) NOT NULL,
  `blob_mode` tinyint(1) NOT NULL,
  `file_path` varchar(255) DEFAULT NULL,
  `file_url` varchar(255) DEFAULT NULL,
  `file_notes` varchar(150) DEFAULT NULL,
  `file_obj` blob DEFAULT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_fund_structure`
--

CREATE TABLE `ms_fund_structure` (
  `fund_structure_key` int(11) NOT NULL,
  `fund_structure_code` varchar(20) NOT NULL,
  `fund_structure_name` varchar(50) NOT NULL,
  `fund_structure_desc` varchar(150) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_fund_structure`
--

INSERT INTO `ms_fund_structure` (`fund_structure_key`, `fund_structure_code`, `fund_structure_name`, `fund_structure_desc`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'OpendEnd', 'Opend_End', 'Opend_End', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'CloseEnd', 'Close_End', 'Close_End', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_fund_type`
--

CREATE TABLE `ms_fund_type` (
  `fund_type_key` int(11) NOT NULL,
  `fund_type_code` varchar(20) DEFAULT NULL,
  `fund_type_name` varchar(50) DEFAULT NULL,
  `fund_type_desc` varchar(150) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_fund_type`
--

INSERT INTO `ms_fund_type` (`fund_type_key`, `fund_type_code`, `fund_type_name`, `fund_type_desc`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'MM', 'Money Market', 'Money Market', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'FI', 'Fix Income', 'Fix Income', 1, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 'BF', 'Balance Fund', 'Balance Fund', 2, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 'EQ', 'Ekuitas', 'Ekuitas', 3, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 'IDX', 'Index', 'Index', 4, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 'PT', 'Penyertaan Terbatas', 'Penyertaan Terbatas', 5, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(7, 'TP', 'Terproteksi', 'Terproteksi', 6, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(8, 'KPD', 'KPD', 'Discre or KPD', 7, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_geolocation`
--

CREATE TABLE `ms_geolocation` (
  `location_key` int(11) NOT NULL,
  `location_code` varchar(20) DEFAULT NULL,
  `location_name` varchar(100) DEFAULT NULL,
  `detect_city` varchar(100) DEFAULT NULL,
  `detect_longitude` varchar(30) DEFAULT NULL,
  `detect_latitude` varchar(30) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_holiday`
--

CREATE TABLE `ms_holiday` (
  `holiday_key` int(11) NOT NULL,
  `stock_market_key` int(11) NOT NULL,
  `holiday_date` date NOT NULL,
  `holiday_name` longtext DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_instrument`
--

CREATE TABLE `ms_instrument` (
  `instrument_key` int(11) NOT NULL,
  `instrument_code` varchar(20) DEFAULT NULL,
  `instrument_name` varchar(50) DEFAULT NULL,
  `instrument_desc` varchar(100) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_instrument`
--

INSERT INTO `ms_instrument` (`instrument_key`, `instrument_code`, `instrument_name`, `instrument_desc`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'CASH', 'Cash and Equivalent', 'Kas', 10, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '1', NULL, NULL),
(2, 'CORPORATE.BOND', 'Corporate Bonds', 'Obligasi Korporasi', 20, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2', NULL, NULL),
(3, 'GOVERNMENT.BOND', 'Government Bonds', 'Obligasi Pemerintah', 30, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '4', NULL, NULL),
(4, 'DEPOSIT', 'Deposito', 'Deposito', 40, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '3', NULL, NULL),
(5, 'EQUITY', 'Equity', 'Saham', 50, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '5', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_participant`
--

CREATE TABLE `ms_participant` (
  `participant_key` int(11) NOT NULL,
  `participant_code` varchar(30) NOT NULL,
  `participant_name` varchar(150) NOT NULL,
  `participant_category` varchar(10) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_participant`
--

INSERT INTO `ms_participant` (`participant_key`, `participant_code`, `participant_name`, `participant_category`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'AAM01', 'PT SHINHAN ASSET MANAGEMENT INDONESIA', 'IM', 1, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'AAM01', 'PT SHINHAN ASSET MANAGEMENT INDONESIA', 'SA', 2, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 'AD001', 'PT OSO SEKURITAS INDONESIA', 'BR', 3, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 'AF001', 'PT HARITA KENCANA SEKURITAS', 'BR', 4, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 'AG001', 'PT KIWOOM SEKURITAS INDONESIA', 'BR', 5, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 'AG002', 'PT KIWOOM INVESTMENT MANAGEMENT INDONESIA', 'IM', 6, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(7, 'AG002', 'PT KIWOOM INVESTMENT MANAGEMENT INDONESIA', 'SA', 7, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(8, 'AH001', 'PT SHINHAN SEKURITAS INDONESIA', 'BR', 8, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(9, 'AH002', 'PT. EMCO ASSET MANAGEMENT', 'IM', 9, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(10, 'AH002', 'PT. EMCO ASSET MANAGEMENT', 'SA', 10, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(11, 'AI001', 'PT UOB KAY HIAN SEKURITAS', 'BR', 11, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(12, 'AK001', 'PT UBS SEKURITAS INDONESIA', 'BR', 12, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(13, 'AMC01', 'PT YUANTA ASSET MANAGEMENT', 'IM', 13, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(14, 'AMC01', 'PT YUANTA ASSET MANAGEMENT', 'SA', 14, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(15, 'AN001', 'PT WANTEG SEKURITAS', 'BR', 15, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(16, 'AN002', 'PT WANTEG ASSET MANAGEMENT', 'IM', 16, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(17, 'AN002', 'PT WANTEG ASSET MANAGEMENT', 'SA', 17, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(18, 'ANA02', 'PT ANARGYA ASET MANAJEMEN', 'IM', 18, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(19, 'ANA02', 'PT ANARGYA ASET MANAJEMEN', 'SA', 19, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(20, 'ANZ05', 'PT BANK ANZ INDONESIA', 'BR', 20, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(21, 'ANZ69', 'PT BANK ANZ INDONESIA', 'SA', 21, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(22, 'AO001', 'ERDIKHA ELIT, PT', 'BR', 22, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(23, 'AP001', 'PT PACIFIC SEKURITAS INDONESIA', 'BR', 23, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(24, 'AP002', 'PT. PACIFIC CAPITAL INVESTMENT', 'IM', 24, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(25, 'AP002', 'PT. PACIFIC CAPITAL INVESTMENT', 'SA', 25, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(26, 'AR001', 'PT BINAARTHA SEKURITAS', 'BR', 26, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(27, 'AR069', 'PT BINAARTHA SEKURITAS', 'SA', 27, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(28, 'ARK02', 'PT ASIA RAYA KAPITAL', 'IM', 28, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(29, 'ARK02', 'PT ASIA RAYA KAPITAL', 'SA', 29, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(30, 'ARO02', 'PT AURORA ASSET MANAGEMENT', 'IM', 30, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(31, 'ARO02', 'PT AURORA ASSET MANAGEMENT', 'SA', 31, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(32, 'ASI02', 'PT NUSADANA INVESTAMA INDONESIA', 'IM', 32, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(33, 'ASI02', 'PT NUSADANA INVESTAMA INDONESIA', 'SA', 33, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(34, 'AT001', 'PT PHINTRACO SEKURITAS', 'BR', 34, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(35, 'AYR02', 'PT AYERS ASIA ASSET MANAGEMENT', 'IM', 35, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(36, 'AYR02', 'PT AYERS ASIA ASSET MANAGEMENT', 'SA', 36, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(37, 'AZ001', 'PT SUCOR SEKURITAS', 'BR', 37, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(38, 'BALI1', 'BANK PERMATA TBK, PT', 'CB', 38, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(39, 'BAM02', 'PT BERLIAN ASET MANAJEMEN', 'IM', 39, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(40, 'BAM02', 'PT BERLIAN ASET MANAJEMEN', 'SA', 40, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(41, 'BATA2', 'PT BATIK ASET MANAJEMEN', 'IM', 41, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(42, 'BATA2', 'PT BATIK ASET MANAJEMEN', 'SA', 42, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(43, 'BBKP2', 'PT BANK BUKOPIN', 'CB', 43, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(44, 'BCA01', 'BANK CENTRAL ASIA TBK, PT', 'CB', 44, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(45, 'BCA05', 'PT.BANK CENTRAL ASIA,TBK', 'BR', 45, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(46, 'BCA69', 'PT BANK CENTRAL ASIA, TBK', 'SA', 46, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(47, 'BCL69', 'PT BUANA CAPITAL', 'SA', 47, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(48, 'BD001', 'PT.INDO MITRA SEKURITAS', 'BR', 48, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(49, 'BDMN2', 'BANK DANAMON INDONESIA TBK, PT', 'CB', 49, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(50, 'BDMN5', 'PT BANK DANAMON INDONESIA, TBK', 'BR', 50, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(51, 'BF001', 'PT. INTI FIKASA SEKURITAS', 'BR', 51, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(52, 'BII01', 'PT BANK MAYBANK INDONESIA TBK', 'CB', 52, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(53, 'BII02', 'PT AXA ASSET MANAGEMENT INDONESIA', 'IM', 53, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(54, 'BII02', 'PT AXA ASSET MANAGEMENT INDONESIA', 'SA', 54, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(55, 'BII05', 'PT BANK MAYBANK INDONESIA TBK', 'BR', 55, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(56, 'BII69', 'PT. BANK MAYBANK INDONESIA, TBK', 'SA', 56, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(57, 'BJ002', 'PT ASANUSA ASSET MANAGEMENT', 'IM', 57, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(58, 'BJ002', 'PT ASANUSA ASSET MANAGEMENT', 'SA', 58, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(59, 'BJB01', 'PT. BANK PEMBANGUNAN DAERAH JAWA BARAT DAN BANTEN TBK', 'CB', 59, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(60, 'BJB69', 'BANK BJB', 'SA', 60, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(61, 'BK001', 'PT J.P. MORGAN SEKURITAS INDONESIA', 'BR', 61, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(62, 'BKP69', 'PT BUKOPIN BANK TBK.', 'SA', 62, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(63, 'BMAN1', 'BANK MANDIRI, PT - CUSTODY', 'CB', 63, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(64, 'BMAN5', 'PT BANK MANDIRI (PERSERO) TBK', 'BR', 64, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(65, 'BMI02', 'PT BUMIPUTERA MANAJEMEN INVESTASI', 'IM', 65, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(66, 'BMI02', 'PT BUMIPUTERA MANAJEMEN INVESTASI', 'SA', 66, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(67, 'BNA69', 'PT. BANK UOB INDONESIA', 'SA', 67, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(68, 'BNGA1', 'BANK CIMB NIAGA TBK, PT', 'CB', 68, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(69, 'BNGA5', 'PT BANK CIMB NIAGA TBK', 'BR', 69, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(70, 'BNI01', 'BANK NEGARA INDONESIA (PERSERO), TBK', 'CB', 70, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(71, 'BNI05', 'PT BANK NEGARA INDONESIA (PERSERO) TBK', 'BR', 71, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(72, 'BNI69', 'PT. BANK NEGARA INDONESIA (PERSERO),TBK', 'SA', 72, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(73, 'BNP05', 'PT BANK BNP PARIBAS INDONESIA', 'BR', 73, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(74, 'BNP69', 'PT BANK BNP PARIBAS INDONESIA', 'SA', 74, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(75, 'BOFA5', 'BANK OF AMERICA NA', 'BR', 75, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(76, 'BPI69', 'PT BAREKSA PORTAL INVESTASI', 'SA', 76, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(77, 'BQ001', 'PT KOREA INVESTMENT AND SEKURITAS INDONESIA', 'BR', 77, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(78, 'BQ002', 'PT KISI ASSET MANAGEMENT', 'IM', 78, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(79, 'BQ002', 'PT KISI ASSET MANAGEMENT', 'SA', 79, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(80, 'BR001', 'PT TRUST SEKURITAS', 'BR', 80, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(81, 'BR069', 'PT TRUST SEKURITAS', 'SA', 81, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(82, 'BRI01', 'BANK RAKYAT INDONESIA (PERSERO), PT', 'CB', 82, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(83, 'BRI05', 'PT BANK RAKYAT INDONESIA (PERSERO) TBK.', 'BR', 83, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(84, 'BRI69', 'PT BANK RAKYAT INDONESIA (PERSERO) TBK', 'SA', 84, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(85, 'BS001', 'PT EQUITY SEKURITAS INDONESIA', 'BR', 85, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(86, 'BS002', 'PT EQUITY SEKURITAS INDONESIA', 'IM', 86, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(87, 'BS002', 'PT EQUITY SEKURITAS INDONESIA', 'SA', 87, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(88, 'BSM01', 'PT BANK SYARIAH MANDIRI', 'CB', 88, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(89, 'BTB69', 'PT BIBIT TUMBUH BERSAMA', 'SA', 89, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(90, 'BTN69', 'PT. BANK TABUNGAN NEGARA (PERSERO) TBK', 'SA', 90, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(91, 'BTP69', 'PT BANK BTPN TBK', 'SA', 91, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(92, 'BWS02', 'PT BOWSPRIT ASSET MANAGEMENT', 'IM', 92, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(93, 'BWS02', 'PT BOWSPRIT ASSET MANAGEMENT', 'SA', 93, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(94, 'BZ001', 'PT. BATAVIA PROSPERINDO SEKURITAS', 'BR', 94, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(95, 'BZ002', 'PT BATAVIA PROSPERINDO ASET MANAJEMEN', 'IM', 95, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(96, 'BZ002', 'PT BATAVIA PROSPERINDO ASET MANAJEMEN', 'SA', 96, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(97, 'CC001', 'MANDIRI SEKURITAS, PT', 'BR', 97, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(98, 'CC002', 'PT MANDIRI MANAJEMEN INVESTASI', 'IM', 98, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(99, 'CC002', 'PT MANDIRI MANAJEMEN INVESTASI', 'SA', 99, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(100, 'CC069', 'PT. MANDIRI SEKURITAS', 'SA', 100, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(101, 'CD001', 'PT MEGA CAPITAL SEKURITAS', 'BR', 101, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(102, 'CD003', 'PT MEGA CAPITAL INVESTAMA', 'IM', 102, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(103, 'CD003', 'PT MEGA CAPITAL INVESTAMA', 'SA', 103, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(104, 'CD004', 'PT MEGA ASSET MANAGEMENT', 'IM', 104, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(105, 'CD004', 'PT MEGA ASSET MANAGEMENT', 'SA', 105, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(106, 'CD069', 'PT MEGA CAPITAL SEKURITAS', 'SA', 106, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(107, 'CFN02', 'PT. CORFINA CAPITAL', 'IM', 107, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(108, 'CFN02', 'PT. CORFINA CAPITAL', 'SA', 108, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(109, 'CG001', 'PT CITIGROUP SEKURITAS INDONESIA', 'BR', 109, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(110, 'CIT69', 'CITIBANK N.A. CABANG INDONESIA', 'SA', 110, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(111, 'CITI1', 'CITIBANK, N. A', 'CB', 111, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(112, 'CITI5', 'CITIBANK. NA JAKARTA', 'BR', 112, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(113, 'COM05', 'PT.BANK COMMONWEALTH', 'BR', 113, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(114, 'COM69', 'PT BANK COMMONWEALTH', 'SA', 114, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(115, 'CP001', 'PT VALBURY SEKURITAS INDONESIA', 'BR', 115, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(116, 'CP002', 'PT VALBURY CAPITAL MANAGEMENT', 'IM', 116, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(117, 'CP002', 'PT VALBURY CAPITAL MANAGEMENT', 'SA', 117, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(118, 'CS002', 'PT CREDIT SUISSE SEKURITAS INDONESIA', 'BR', 118, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(119, 'CTB69', 'BANK CTBC INDONESIA', 'SA', 119, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(120, 'DB001', 'PT DEUTSCHE SEKURITAS INDONESIA', 'BR', 120, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(121, 'DBJK1', 'BUT DEUTSCHE BANK AG', 'CB', 121, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(122, 'DBJK5', 'DEUTSCHE BANK AG', 'BR', 122, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(123, 'DBS69', 'BANK DBS INDONESIA', 'SA', 123, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(124, 'DBSI1', 'PT BANK DBS INDONESIA', 'CB', 124, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(125, 'DBSI5', 'PT BANK DBS INDONESIA', 'BR', 125, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(126, 'DD001', 'PT MAKINDO SEKURITAS', 'BR', 126, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(127, 'DD002', 'INTRU NUSANTARA, PT', 'IM', 127, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(128, 'DD002', 'INTRU NUSANTARA, PT', 'SA', 128, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(129, 'DD069', 'PT MAKINDO SECURITIES', 'SA', 129, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(130, 'DH001', 'SINARMAS SEKURITAS, PT', 'BR', 130, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(131, 'DH002', 'PT SINARMAS ASSET MANAGEMENT', 'IM', 131, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(132, 'DH002', 'PT SINARMAS ASSET MANAGEMENT', 'SA', 132, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(133, 'DH069', 'PT. SINARMAS SEKURITAS', 'SA', 133, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(134, 'DI002', 'PT DANAKITA INVESTAMA', 'IM', 134, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(135, 'DI002', 'PT DANAKITA INVESTAMA', 'SA', 135, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(136, 'DMN69', 'PT. BANK DANAMON INDONESIA, TBK', 'SA', 136, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(137, 'DP001', 'PT DBS VICKERS SEKURITAS INDONESIA', 'BR', 137, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(138, 'DR001', 'PT RHB SEKURITAS INDONESIA', 'BR', 138, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(139, 'DR002', 'PT RHB ASSET MANAGEMENT INDONESIA', 'IM', 139, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(140, 'DR002', 'PT RHB ASSET MANAGEMENT INDONESIA', 'SA', 140, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(141, 'DR069', 'RHB SECURITIES INDONESIA, PT', 'SA', 141, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(142, 'DX001', 'PT BAHANA SEKURITAS', 'BR', 142, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(143, 'DX002', 'PT. BAHANA TCW INVESTMENT MANAGEMENT', 'IM', 143, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(144, 'DX002', 'PT. BAHANA TCW INVESTMENT MANAGEMENT', 'SA', 144, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(145, 'EII01', 'PT. EASTSPRING INVESTMENTS INDONESIA', 'IM', 145, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(146, 'EII01', 'PT. EASTSPRING INVESTMENTS INDONESIA', 'SA', 146, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(147, 'EL001', 'PT. EVERGREEN SEKURITAS INDONESIA', 'BR', 147, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(148, 'EL002', 'PT ASHMORE ASSET MANAGEMENT INDONESIA', 'IM', 148, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(149, 'EL002', 'PT ASHMORE ASSET MANAGEMENT INDONESIA', 'SA', 149, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(150, 'EP001', 'PT MNC SEKURITAS', 'BR', 150, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(151, 'EP002', 'PT MNC ASSET MANAGEMENT', 'IM', 151, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(152, 'EP002', 'PT MNC ASSET MANAGEMENT', 'SA', 152, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(153, 'ES001', 'EKOKAPITAL SEKURITAS, PT', 'BR', 153, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(154, 'FG001', 'PT NOMURA SEKURITAS INDONESIA', 'BR', 154, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(155, 'FM001', 'PT ONIX SEKURITAS', 'BR', 155, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(156, 'FS001', 'PT YUANTA SEKURITAS INDONESIA', 'BR', 156, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(157, 'FSI01', 'PT FIRST STATE INVESTMENTS INDONESIA', 'IM', 157, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(158, 'FSI01', 'PT FIRST STATE INVESTMENTS INDONESIA', 'SA', 158, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(159, 'FZ001', 'PT WATERFRONT SEKURITAS INDONESIA', 'BR', 159, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(160, 'GA001', 'PT BNC SEKURITAS INDONESIA', 'BR', 160, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(161, 'GAMA2', 'PT SUCORINVEST ASSET MANAGEMENT', 'IM', 161, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(162, 'GAMA2', 'PT SUCORINVEST ASSET MANAGEMENT', 'SA', 162, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(163, 'GAP01', 'PT GAP CAPITAL', 'IM', 163, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(164, 'GAP01', 'PT GAP CAPITAL', 'SA', 164, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(165, 'GMI02', 'PT. GEMILANG INDONESIA MANAJEMEN INVESTASI', 'IM', 165, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(166, 'GMI02', 'PT. GEMILANG INDONESIA MANAJEMEN INVESTASI', 'SA', 166, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(167, 'GMT01', 'PT MAYBANK ASSET MANAGEMENT', 'IM', 167, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(168, 'GMT01', 'PT MAYBANK ASSET MANAGEMENT', 'SA', 168, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(169, 'GNES5', 'PT BANK GANESHA TBK.', 'BR', 169, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(170, 'GNS69', 'PT. BANK GANESHA TBK', 'SA', 170, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(171, 'GR001', 'PANIN SEKURITAS TBK, PT', 'BR', 171, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(172, 'GR003', 'PT PANIN ASSET MANAGEMENT', 'IM', 172, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(173, 'GR003', 'PT PANIN ASSET MANAGEMENT', 'SA', 173, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(174, 'GW001', 'PT HSBC SEKURITAS INDONESIA', 'BR', 174, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(175, 'HD001', 'PT. KGI SEKURITAS INDONESIA', 'BR', 175, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(176, 'HK002', 'PT OSO MANAJEMEN INVESTASI', 'IM', 176, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(177, 'HK002', 'PT OSO MANAJEMEN INVESTASI', 'SA', 177, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(178, 'HNA69', 'PT BANK KEB HANA INDONESIA', 'SA', 178, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(179, 'HP001', 'PT HENAN PUTIHRAI SEKURITAS', 'BR', 179, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(180, 'HP002', 'PT HENAN PUTIHRAI ASSET MANAGEMENT', 'IM', 180, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(181, 'HP002', 'PT HENAN PUTIHRAI ASSET MANAGEMENT', 'SA', 181, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(182, 'HPS69', 'PT HENAN PUTIHRAI', 'SA', 182, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(183, 'HSB69', 'PT BANK HSBC INDONESIA', 'SA', 183, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(184, 'HSBC1', 'PT BANK HSBC INDONESIA', 'CB', 184, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(185, 'HSBC5', 'PT BANK HSBC INDONESIA', 'BR', 185, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(186, 'IAI02', 'PT. INDO ARTHABUANA INVESTAMA', 'IM', 186, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(187, 'IAI02', 'PT. INDO ARTHABUANA INVESTAMA', 'SA', 187, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(188, 'IAM02', 'PT INDOSTERLING ASET MANAJEMEN', 'IM', 188, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(189, 'IAM02', 'PT INDOSTERLING ASET MANAJEMEN', 'SA', 189, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(190, 'ID001', 'PT ANUGERAH SEKURITAS INDONESIA', 'BR', 190, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(191, 'ID002', 'PT ANUGERAH SENTRA INVESTAMA', 'IM', 191, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(192, 'ID002', 'PT ANUGERAH SENTRA INVESTAMA', 'SA', 192, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(193, 'IF001', 'SAMUEL SEKURITAS INDONESIA, PT', 'BR', 193, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(194, 'IF002', 'PT SAMUEL ASET MANAJEMEN', 'IM', 194, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(195, 'IF002', 'PT SAMUEL ASET MANAJEMEN', 'SA', 195, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(196, 'II001', 'PT DANATAMA MAKMUR SEKURITAS', 'BR', 196, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(197, 'II002', 'PT DANATAMA MAKMUR', 'IM', 197, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(198, 'II002', 'PT DANATAMA MAKMUR', 'SA', 198, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(199, 'IIM02', 'PT INSIGHT INVESTMENTS MANAGEMENT', 'IM', 199, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(200, 'IIM02', 'PT INSIGHT INVESTMENTS MANAGEMENT', 'SA', 200, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(201, 'IN001', 'INVESTINDO NUSANTARA SEKURITAS, PT', 'BR', 201, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(202, 'IN002', 'PT NARADA ASET MANAJEMEN', 'IM', 202, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(203, 'IN002', 'PT NARADA ASET MANAJEMEN', 'SA', 203, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(204, 'IND02', 'PT INDOASIA ASET MANAJEMEN', 'IM', 204, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(205, 'IND02', 'PT INDOASIA ASET MANAJEMEN', 'SA', 205, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(206, 'INO69', 'PT INDO CAPITAL SEKURITAS', 'SA', 206, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(207, 'INV69', 'PT INVESTAMART PRINCIPAL OPTIMA', 'SA', 207, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(208, 'IP001', 'PT INDOSURYA BERSINAR SEKURITAS', 'BR', 208, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(209, 'IP002', 'PT. INDOSURYA ASSET MANAGEMENT', 'IM', 209, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(210, 'IP002', 'PT. INDOSURYA ASSET MANAGEMENT', 'SA', 210, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(211, 'IPS69', 'PT INDO PREMIER SEKURITAS', 'SA', 211, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(212, 'IU001', 'PT INDO CAPITAL SEKURITAS', 'BR', 212, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(213, 'JAM02', 'PT JARVIS ASET MANAJEMEN', 'IM', 213, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(214, 'JAM02', 'PT JARVIS ASET MANAJEMEN', 'SA', 214, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(215, 'JFI02', 'PT CORPUS KAPITAL MANAJEMEN', 'IM', 215, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(216, 'JFI02', 'PT CORPUS KAPITAL MANAJEMEN', 'SA', 216, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(217, 'JPM05', 'JP MORGAN CHASE, N.A, JAKARTA BRANCH', 'BR', 217, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(218, 'KHANA', 'PT KEB HANA BANK INDONESIA', 'CB', 218, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(219, 'KI001', 'PT CIPTADANA SEKURITAS ASIA', 'BR', 219, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(220, 'KI002', 'PT CIPTADANA ASSET MANAGEMENT', 'IM', 220, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(221, 'KI002', 'PT CIPTADANA ASSET MANAGEMENT', 'SA', 221, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(222, 'KI069', 'PT CIPTADANA SEKURITAS ASIA', 'SA', 222, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(223, 'KK001', 'PT PHILLIP SEKURITAS INDONESIA', 'BR', 223, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(224, 'KK002', 'PT. PHILLIP ASSET MANAGEMENT', 'IM', 224, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(225, 'KK002', 'PT. PHILLIP ASSET MANAGEMENT', 'SA', 225, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(226, 'KK069', 'PT PHILLIP SEKURITAS INDONESIA', 'SA', 226, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(227, 'KS001', 'PT KRESNA SEKURITAS', 'BR', 227, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(228, 'KS002', 'PT KRESNA ASSET MANAGEMENT', 'IM', 228, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(229, 'KS002', 'PT KRESNA ASSET MANAGEMENT', 'SA', 229, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(230, 'KW001', 'PT. CORPUS SEKURITAS INDONESIA', 'BR', 230, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(231, 'KZ001', 'PT CLSA SEKURITAS INDONESIA', 'BR', 231, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(232, 'LG001', 'PT. TRIMEGAH SEKURITAS INDONESIA TBK', 'BR', 232, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(233, 'LG002', 'PT TRIMEGAH ASSET MANAGEMENT', 'IM', 233, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(234, 'LG002', 'PT TRIMEGAH ASSET MANAGEMENT', 'SA', 234, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(235, 'LH001', 'PT ROYAL INVESTIUM SEKURITAS', 'BR', 235, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(236, 'LK002', 'PT RECAPITAL ASSET MANAGEMENT', 'IM', 236, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(237, 'LK002', 'PT RECAPITAL ASSET MANAGEMENT', 'SA', 237, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(238, 'LPPS2', 'PT LIPPO SECURITIES, TBK.', 'IM', 238, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(239, 'LPPS2', 'PT LIPPO SECURITIES, TBK.', 'SA', 239, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(240, 'LS001', 'PT RELIANCE SEKURITAS INDONESIA, TBK', 'BR', 240, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(241, 'LS002', 'PT RELIANCE MANAJER INVESTASI', 'IM', 241, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(242, 'LS002', 'PT RELIANCE MANAJER INVESTASI', 'SA', 242, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(243, 'LS069', 'PT RELIANCE SEKURITAS INDONESIA, TBK', 'SA', 243, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(244, 'MAM02', 'PT MASERI ASET MANAJEMEN', 'IM', 244, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(245, 'MAM02', 'PT MASERI ASET MANAJEMEN', 'SA', 245, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(246, 'MAN02', 'PT MANULIFE ASET MANAJEMEN INDONESIA', 'IM', 246, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(247, 'MAN02', 'PT MANULIFE ASET MANAJEMEN INDONESIA', 'SA', 247, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(248, 'MAN69', 'PT BANK MANDIRI (PERSERO) TBK.', 'SA', 248, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(249, 'MDI69', 'PT. MODUIT DIGITAL INDONESIA', 'SA', 249, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(250, 'MEES2', 'PT. BNP PARIBAS ASSET MANAGEMENT', 'IM', 250, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(251, 'MEES2', 'PT. BNP PARIBAS ASSET MANAGEMENT', 'SA', 251, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(252, 'MEGA3', 'BANK MEGA TBK, PT', 'CB', 252, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(253, 'MEGA5', 'PT. BANK MEGA TBK.', 'BR', 253, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(254, 'MG001', 'PT SEMESTA INDOVEST SEKURITAS', 'BR', 254, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(255, 'MG002', 'PT SEMESTA ASET MANAJEMEN', 'IM', 255, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(256, 'MG002', 'PT SEMESTA ASET MANAJEMEN', 'SA', 256, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(257, 'MGA69', 'PT BANK MEGA, TBK.', 'SA', 257, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(258, 'MI001', 'PT VICTORIA SEKURITAS INDONESIA', 'BR', 258, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(259, 'MJR02', 'PT MAJORIS', 'IM', 259, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(260, 'MJR02', 'PT MAJORIS', 'SA', 260, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(261, 'MK001', 'PT EKUATOR SWARNA SEKURITAS', 'BR', 261, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(262, 'MK002', 'PT EKUATOR SWARNA INVESTAMA', 'IM', 262, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(263, 'MK002', 'PT EKUATOR SWARNA INVESTAMA', 'SA', 263, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(264, 'ML001', 'PT MERRILL LYNCH SEKURITAS INDONESIA', 'BR', 264, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(265, 'MNC69', 'PT. MNC SEKURITAS', 'SA', 265, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(266, 'MS002', 'PT. MORGAN STANLEY SEKURITAS INDONESIA', 'BR', 266, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(267, 'MU001', 'PT MINNA PADI INVESTAMA SEKURITAS TBK', 'BR', 267, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(268, 'MU002', 'PT MINNA PADI ASET MANAJEMEN', 'IM', 268, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(269, 'MU002', 'PT MINNA PADI ASET MANAJEMEN', 'SA', 269, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(270, 'MYR69', 'PT BANK MAYORA', 'SA', 270, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(271, 'NAR02', 'PT CAPITAL ASSET MANAGEMENT', 'IM', 271, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(272, 'NAR02', 'PT CAPITAL ASSET MANAGEMENT', 'SA', 272, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(273, 'NGA69', 'PT. BANK CIMB NIAGA', 'SA', 273, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(274, 'NI001', 'BNI SEKURITAS, PT', 'BR', 274, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(275, 'NI002', 'PT BNI ASSET MANAGEMENT', 'IM', 275, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(276, 'NI002', 'PT BNI ASSET MANAGEMENT', 'SA', 276, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(277, 'NI069', 'PT BNI SECURITIES', 'SA', 277, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(278, 'NIB69', 'PT NADIRA INVESTASIKITA BERSAMA', 'SA', 278, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(279, 'NOB69', 'PT. BANK NATIONALNOBU TBK.', 'SA', 279, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(280, 'NSI69', 'PT. NUSANTARA SEJAHTERA INVESTAMA', 'SA', 280, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(281, 'NSP05', 'PT. BANK OCBC NISP TBK', 'BR', 281, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(282, 'NSP69', 'PT. BANK OCBC NISP TBK', 'SA', 282, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(283, 'NUS02', 'PT NUSANTARA SENTRA KAPITAL', 'IM', 283, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(284, 'NUS02', 'PT NUSANTARA SENTRA KAPITAL', 'SA', 284, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(285, 'OD001', 'DANAREKSA SEKURITAS, PT', 'BR', 285, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(286, 'OD002', 'PT DANAREKSA INVESTMENT MANAGEMENT', 'IM', 286, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(287, 'OD002', 'PT DANAREKSA INVESTMENT MANAGEMENT', 'SA', 287, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(288, 'OD069', 'PT DANAREKSA SEKURITAS', 'SA', 288, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(289, 'OK001', 'NET SEKURITAS, PT', 'BR', 289, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(290, 'OK002', 'PT NET ASSETS MANAGEMENT', 'IM', 290, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(291, 'OK002', 'PT NET ASSETS MANAGEMENT', 'SA', 291, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(292, 'PAMI2', 'PT POST ASSET MANAGEMENT INDONESIA', 'IM', 292, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(293, 'PAMI2', 'PT POST ASSET MANAGEMENT INDONESIA', 'SA', 293, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(294, 'PAN05', 'PT. BANK PAN INDONESIA TBK', 'BR', 294, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(295, 'PAN69', 'PT. BANK PAN INDONESIA TBK.', 'SA', 295, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(296, 'PD001', 'PT INDO PREMIER SEKURITAS', 'BR', 296, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(297, 'PD002', 'PT INDO PREMIER INVESTMENT MANAGEMENT', 'IM', 297, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(298, 'PD002', 'PT INDO PREMIER INVESTMENT MANAGEMENT', 'SA', 298, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(299, 'PG001', 'PT PANCA GLOBAL SEKURITAS', 'BR', 299, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(300, 'PG003', 'PT. UOB ASSET MANAGEMENT INDONESIA', 'IM', 300, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(301, 'PG003', 'PT. UOB ASSET MANAGEMENT INDONESIA', 'SA', 301, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(302, 'PI001', 'PT.MAGENTA KAPITAL SEKURITAS INDONESIA', 'BR', 302, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(303, 'PI002', 'PT PINNACLE PERSADA INVESTAMA', 'IM', 303, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(304, 'PI002', 'PT PINNACLE PERSADA INVESTAMA', 'SA', 304, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(305, 'PK001', 'PT PRATAMA CAPITAL SEKURITAS', 'BR', 305, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(306, 'PLA02', 'PT PRATAMA CAPITAL ASSETS MANAGEMENT', 'IM', 306, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(307, 'PLA02', 'PT PRATAMA CAPITAL ASSETS MANAGEMENT', 'SA', 307, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(308, 'PMT05', 'PT. BANK PERMATA TBK.', 'BR', 308, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(309, 'PMT69', 'BANK PERMATA TBK, PT', 'SA', 309, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(310, 'PNM02', 'PT PNM INVESTMENT MANAGEMENT', 'IM', 310, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(311, 'PNM02', 'PT PNM INVESTMENT MANAGEMENT', 'SA', 311, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(312, 'PO001', 'PT PILARMAS INVESTINDO SEKURITAS', 'BR', 312, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(313, 'PP001', 'PT ALDIRACITA SEKURITAS INDONESIA', 'BR', 313, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(314, 'PP069', 'PT ALDIRACITA SEKURITAS INDONESIA', 'SA', 314, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(315, 'PRM02', 'PT. JASA CAPITAL ASSET MANAGEMENT', 'IM', 315, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(316, 'PRM02', 'PT. JASA CAPITAL ASSET MANAGEMENT', 'SA', 316, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(317, 'PROS2', 'PT.PROSPERA ASSET MANAGEMENT', 'IM', 317, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(318, 'PROS2', 'PT.PROSPERA ASSET MANAGEMENT', 'SA', 318, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(319, 'PS001', 'PARAMITRA ALFA SEKURITAS, PT', 'BR', 319, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(320, 'PS002', 'PT. PARAMITRA ALFA SEKURITAS', 'IM', 320, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(321, 'PS002', 'PT. PARAMITRA ALFA SEKURITAS', 'SA', 321, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(322, 'PTM02', 'PT PAYTREN ASET MANAJEMEN', 'IM', 322, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(323, 'PTM02', 'PT PAYTREN ASET MANAJEMEN', 'SA', 323, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(324, 'QA001', 'PT POOL ADVISTA SEKURITAS', 'BR', 324, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(325, 'QK002', 'PT FOSTER ASSET MANAGEMENT', 'IM', 325, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(326, 'QK002', 'PT FOSTER ASSET MANAGEMENT', 'SA', 326, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(327, 'QNB69', 'BANK QNB INDONESIA', 'SA', 327, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(328, 'RA002', 'PT POOL ADVISTA ASET MANAJEMEN', 'IM', 328, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(329, 'RA002', 'PT POOL ADVISTA ASET MANAJEMEN', 'SA', 329, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(330, 'RAI69', 'PT RAIZ INVEST INDONESIA', 'SA', 330, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(331, 'RAM02', 'PT RAHA ASET MANAJEMEN', 'IM', 331, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(332, 'RAM02', 'PT RAHA ASET MANAJEMEN', 'SA', 332, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(333, 'RB001', 'PT NIKKO SEKURITAS INDONESIA', 'BR', 333, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO `ms_participant` (`participant_key`, `participant_code`, `participant_name`, `participant_category`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(334, 'RB002', 'PT NIKKO SEKURITAS INDONESIA', 'IM', 334, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(335, 'RB002', 'PT NIKKO SEKURITAS INDONESIA', 'SA', 335, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(336, 'RF001', 'PT BUANA CAPITAL SEKURITAS', 'BR', 336, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(337, 'RG001', 'PT PROFINDO SEKURITAS INDONESIA', 'BR', 337, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(338, 'RO001', 'NISP SEKURITAS, PT', 'BR', 338, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(339, 'RO002', 'PT ABERDEEN STANDARD INVESTMENTS INDONESIA', 'IM', 339, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(340, 'RO002', 'PT ABERDEEN STANDARD INVESTMENTS INDONESIA', 'SA', 340, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(341, 'RX001', 'PT MACQUARIE SEKURITAS INDONESIA', 'BR', 341, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(342, 'SA001', 'PT BOSOWA SEKURITAS', 'BR', 342, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(343, 'SC001', 'PT IMG SEKURITAS', 'BR', 343, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(344, 'SCB05', 'STANDARD CHARTERED BANK', 'BR', 344, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(345, 'SCB69', 'STANDARD CHARTERED BANK', 'SA', 345, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(346, 'SCBJK', 'BUT. STANDARD CHARTERED BANK', 'CB', 346, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(347, 'SCG69', 'PT SUCOR SEKURITAS', 'SA', 347, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(348, 'SCH02', 'PT SCHRODER INVESTMENT MANAGEMENT INDONESIA', 'IM', 348, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(349, 'SCH02', 'PT SCHRODER INVESTMENT MANAGEMENT INDONESIA', 'SA', 349, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(350, 'SF001', 'PT SURYA FAJAR SEKURITAS', 'BR', 350, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(351, 'SH001', 'PT ARTHA SEKURITAS INDONESIA', 'BR', 351, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(352, 'SHN02', 'PT SHINOKEN ASSET MANAGEMENT INDONESIA', 'IM', 352, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(353, 'SHN02', 'PT SHINOKEN ASSET MANAGEMENT INDONESIA', 'SA', 353, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(354, 'SIM02', 'PT. SETIABUDI INVESTMENT MANAGEMENT', 'IM', 354, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(355, 'SIM02', 'PT. SETIABUDI INVESTMENT MANAGEMENT', 'SA', 355, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(356, 'SM002', 'PT MILLENIUM CAPITAL MANAGEMENT', 'IM', 356, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(357, 'SM002', 'PT MILLENIUM CAPITAL MANAGEMENT', 'SA', 357, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(358, 'SMC69', 'PT STAR MERCATO CAPITALE', 'SA', 358, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(359, 'SNM69', 'BANK SINARMAS', 'SA', 359, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(360, 'SQ001', 'PT BCA SEKURITAS', 'BR', 360, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(361, 'SQAM2', 'PT SEQUIS ASET MANAJEMEN', 'IM', 361, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(362, 'SQAM2', 'PT SEQUIS ASET MANAJEMEN', 'SA', 362, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(363, 'SRI69', 'PT SUPERMARKET REKSADANA INDONESIA', 'SA', 363, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(364, 'STA01', 'PT. SURYA TIMUR ALAM RAYA', 'IM', 364, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(365, 'STA01', 'PT. SURYA TIMUR ALAM RAYA', 'SA', 365, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(366, 'SYA02', 'PT SYAILENDRA CAPITAL', 'IM', 366, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(367, 'SYA02', 'PT SYAILENDRA CAPITAL', 'SA', 367, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(368, 'SYM69', 'PT BANK SYARIAH MANDIRI', 'SA', 368, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(369, 'TF002', 'PT TREASURE FUND INVESTAMA', 'IM', 369, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(370, 'TF002', 'PT TREASURE FUND INVESTAMA', 'SA', 370, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(371, 'TMS69', 'PT TRIMEGAH SEKURITAS INDONESIA TBK', 'SA', 371, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(372, 'TP001', 'PT OCBC SEKURITAS INDONESIA', 'BR', 372, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(373, 'TP002', 'PT AVRIST ASSET MANAGEMENT', 'IM', 373, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(374, 'TP002', 'PT AVRIST ASSET MANAGEMENT', 'SA', 374, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(375, 'TTI69', 'PT TAKJUB TEKNOLOGI INDONESIA', 'SA', 375, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(376, 'TX001', 'PT DHANAWIBAWA SEKURITAS INDONESIA', 'BR', 376, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(377, 'TX002', 'PT PAN ARCADIA CAPITAL', 'IM', 377, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(378, 'TX002', 'PT PAN ARCADIA CAPITAL', 'SA', 378, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(379, 'UOBB5', 'PT. BANK UOB INDONESIA', 'BR', 379, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(380, 'VIC69', 'PT. BANK VICTORIA INTERNATIONAL, TBK', 'SA', 380, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(381, 'VMI02', 'PT VICTORIA MANAJEMEN INVESTASI', 'IM', 381, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(382, 'VMI02', 'PT VICTORIA MANAJEMEN INVESTASI', 'SA', 382, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(383, 'WII69', 'PT WAHED INVESTASI INDONESIA', 'SA', 383, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(384, 'XA001', 'PT NH KORINDO SEKURITAS INDONESIA', 'BR', 384, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(385, 'XID69', 'PT XDANA INVESTA INDONESIA', 'SA', 385, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(386, 'XL001', 'PT MAHAKARYA ARTHA SEKURITAS', 'BR', 386, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(387, 'YB001', 'PT JASA UTAMA CAPITAL SEKURITAS', 'BR', 387, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(388, 'YJ001', 'PT LOTUS ANDALAN SEKURITAS', 'BR', 388, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(389, 'YJ002', 'PT. LAUTANDHANA INVESTMENT MANAGEMENT', 'IM', 389, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(390, 'YJ002', 'PT. LAUTANDHANA INVESTMENT MANAGEMENT', 'SA', 390, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(391, 'YO001', 'PT AMANTARA SEKURITAS INDONESIA', 'BR', 391, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(392, 'YP001', 'PT MIRAE ASSET SEKURITAS INDONESIA', 'BR', 392, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(393, 'YP069', 'PT. MIRAE ASSET SEKURITAS INDONESIA', 'SA', 393, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(394, 'YU001', 'PT CGS-CIMB SEKURITAS INDONESIA', 'BR', 394, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(395, 'Z0001', 'PT PEAK SEKURITAS INDONESIA ', 'BR', 395, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(396, 'ZG002', 'PT PRINCIPAL ASSET MANAGEMENT', 'IM', 396, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(397, 'ZG002', 'PT PRINCIPAL ASSET MANAGEMENT', 'SA', 397, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(398, 'ZP001', 'MAYBANK KIM ENG SECURITIES, PT', 'BR', 398, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(399, 'ZP069', 'PT MAYBANK KIM ENG SEKURITAS', 'SA', 399, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(400, 'ZR001', 'PT BUMIPUTERA SEKURITAS', 'BR', 400, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_product`
--

CREATE TABLE `ms_product` (
  `product_key` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `product_code` varchar(30) NOT NULL,
  `product_name` varchar(150) NOT NULL,
  `product_name_alt` varchar(150) NOT NULL,
  `currency_key` int(11) DEFAULT NULL,
  `product_category_key` int(11) DEFAULT NULL,
  `product_type_key` int(11) DEFAULT NULL,
  `fund_type_key` int(11) DEFAULT NULL,
  `fund_structure_key` int(11) DEFAULT NULL,
  `risk_profile_key` int(11) DEFAULT NULL,
  `product_profile` varchar(500) DEFAULT NULL,
  `investment_objectives` varchar(500) DEFAULT NULL,
  `logo_link` varchar(255) DEFAULT NULL,
  `banner_link` varchar(255) DEFAULT NULL,
  `prospectus_link` varchar(255) DEFAULT NULL,
  `launch_date` date DEFAULT NULL,
  `inception_date` date DEFAULT NULL,
  `isin_code` varchar(50) DEFAULT NULL,
  `flag_syariah` tinyint(1) NOT NULL,
  `max_sub_fee` decimal(9,2) NOT NULL,
  `max_red_fee` decimal(9,2) NOT NULL,
  `max_swi_fee` decimal(9,2) NOT NULL,
  `min_sub_amount` decimal(9,2) NOT NULL,
  `min_red_amount` decimal(9,2) NOT NULL,
  `min_red_unit` decimal(9,2) NOT NULL,
  `min_unit_after_red` decimal(9,2) NOT NULL,
  `management_fee` decimal(8,4) NOT NULL,
  `custodian_fee` decimal(8,4) NOT NULL,
  `custodian_key` int(11) DEFAULT NULL,
  `ojk_fee` decimal(8,4) NOT NULL,
  `product_fee_amount` decimal(9,2) NOT NULL,
  `product_fee_date_start` datetime(6) DEFAULT NULL,
  `product_fee_date_thru` datetime(6) DEFAULT NULL,
  `other_fee_amount` decimal(5,2) NOT NULL,
  `settlement_period` int(11) DEFAULT NULL,
  `sinvest_fund_code` varchar(30) DEFAULT NULL,
  `flag_enabled` tinyint(1) NOT NULL,
  `flag_subscription` tinyint(1) NOT NULL,
  `flag_redemption` tinyint(1) NOT NULL,
  `flag_switch_out` tinyint(1) NOT NULL,
  `flag_switch_in` tinyint(1) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_product`
--

INSERT INTO `ms_product` (`product_key`, `product_id`, `product_code`, `product_name`, `product_name_alt`, `currency_key`, `product_category_key`, `product_type_key`, `fund_type_key`, `fund_structure_key`, `risk_profile_key`, `product_profile`, `investment_objectives`, `logo_link`, `banner_link`, `prospectus_link`, `launch_date`, `inception_date`, `isin_code`, `flag_syariah`, `max_sub_fee`, `max_red_fee`, `max_swi_fee`, `min_sub_amount`, `min_red_amount`, `min_red_unit`, `min_unit_after_red`, `management_fee`, `custodian_fee`, `custodian_key`, `ojk_fee`, `product_fee_amount`, `product_fee_date_start`, `product_fee_date_thru`, `other_fee_amount`, `settlement_period`, `sinvest_fund_code`, `flag_enabled`, `flag_subscription`, `flag_redemption`, `flag_switch_out`, `flag_switch_in`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 5, 'BDLC', 'MNC DANA LANCAR', 'DANA LANCAR', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MMCDANLCR00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 11, 'BDM', 'MNC DANA SYARIAH', 'DANA SYARIAH', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FISDANSYA00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 12, 'BDLSATU', 'MNC DANA LIKUID', 'DANA LIKUID', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FICDANLIK00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 14, 'BBK', 'MNC DANA KOMBINASI', 'DANA KOMBINASI', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MXCDANKOM00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 16, 'BBE', 'MNC DANA EKUITAS', 'DANA EKUITAS', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002EQCDANEKT00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 186, 'MNCUSD', 'MNC DANA DOLLAR', 'DANA DOLLAR', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FICDANDLR00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(7, 187, 'ICON', 'MNC DANA KOMBINASI ICON', 'DANA KOMBINASI ICON', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MXCDANKOI00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(8, 188, 'MDSK', 'MNC DANA SYARIAH KOMBINASI', 'DANA SYARIAH KOMBINASI', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MXSDANKOS00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(9, 189, 'MDSE', 'MNC DANA SYARIAH EKUITAS', 'DANA SYARIAH EKUITAS', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002EQSDANSEK00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(10, 190, 'MDKK', 'MNC DANA KOMBINASI KONSUMEN', 'DANA KOMBINASI KONSUMEN', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MXCDANKOK00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(11, 193, 'MNCLP', 'KPD MNC LINK PASTI', 'MNC LINK PASTI', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002DFCKPDLPA00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(12, 194, 'MNCLA', 'KPD MNC LINK AKTIF', 'MNC LINK AKTIF', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002DFCKPDLAK00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(13, 195, 'MNCLS', 'KPD MNC LINK SERASI', 'MNC LINK SERASI', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002DFCKPDLSE00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(14, 197, 'MNCDTIII', 'MNC DANA TERPROTEKSI III', 'DANA TERPROTEKSI III', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO03', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(15, 198, 'MNCPTII', 'MNC DANA PENDAPATAN TETAP II', 'DANA PENDAPATAN TETAP II', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FICDANTTP02', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(16, 199, 'MNC36', 'REKSA DANA INDEKS MNC36', 'INDEKS MNC36', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002IFCIDKS3600', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(17, 201, 'MNCDTV', 'MNC DANA TERPROTEKSI V', 'DANA TERPROTEKSI V', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO05', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(18, 204, 'MNCSBAR', 'REKSA DANA SYARIAH MNC DANA SYARIAH BAROKAH', 'DANA SYARIAH BAROKAH', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MMSDANSYB00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(19, 205, 'MNCDTVIII', 'MNC DANA TERPROTEKSI VIII', 'DANA TERPROTEKSI VIII', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO08', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(20, 206, 'MNCSBN', 'REKSA DANA MNC DANA SBN', 'DANA SBN', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FICDANSBN00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(21, 207, 'MNCPTIII', 'MNC DANA PENDAPATAN TETAP III', 'DANA PENDAPATAN TETAP III', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FICDANTTP03', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(22, 208, 'MNCSEF', 'MNC SMART EQUITY FUND', 'SMART EQUITY FUND', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002EQCDANSHM02', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(23, 209, 'MNCDTIX', 'MNC DANA TERPROTEKSI IX', 'DANA TERPROTEKSI IX', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO09', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(24, 210, 'MNCDTXII', 'MNC DANA TERPROTEKSI XII', 'DANA TERPROTEKSI XII', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO12', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(25, 211, 'MNCDTX', 'MNC DANA TERPROTEKSI X', 'DANA TERPROTEKSI X', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCMNCDTX10', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(26, 213, 'MDSEII', 'MNC DANA SYARIAH EKUITAS II', 'DANA SYARIAH EKUITAS II', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002EQSDANSEK02', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(27, 215, 'MNCDTXIII', 'MNC DANA TERPROTEKSI XIII', 'DANA TERPROTEKSI XIII', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO13', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(28, 216, 'MNCSEKTOR1', 'REKSADANA PENYERTAAN TERBATAS MNC DANA MULTISEKTOR I', 'DANA MULTISEKTOR I', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PECSEKTOR01', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(29, 217, 'MNCDTXIV', 'MNC DANA TERPROTEKSI XIV', 'DANA TERPROTEKSI XIV', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO14', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(30, 218, 'MNCPTIV', 'MNC DANA PENDAPATAN TETAP IV', 'DANA PENDAPATAN TETAP IV', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FICDANTTP04', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(31, 220, 'MNCDTXVII', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI XVII', 'DANA TERPROTEKSI XVII', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO17', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(32, 221, 'MNCDPUII', 'REKSA DANA MNC DANA PASAR UANG II', 'DANA PASAR UANG II', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MMCRDMDPU02', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(33, 223, 'MNCDTXVIII', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI XVIII', 'DANA TERPROTEKSI XVIII', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO18', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(34, 224, 'MNCDTXIX', 'REKSADANA TERPROTEKSI MNC DANA TERPROTEKSI XIX', 'DANA TERPROTEKSI XIX', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO19', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(35, 225, 'MNCDPU', 'REKSA DANA MNC DANA PASAR UANG', 'DANA PASAR UANG', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002MMCRDMDPU00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(36, 226, 'MNCPTV', 'REKSA DANA MNC DANA PENDAPATAN TETAP V', 'DANA PENDAPATAN TETAP V', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FICDANTTP05', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(37, 227, 'MNCDT21', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 21', 'DANA TERPROTEKSI SERI 21', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO21', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(38, 228, 'MNCDT23', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 23', 'DANA TERPROTEKSI SERI 23', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO23', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(39, 229, 'MNCDT22', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 22', 'DANA TERPROTEKSI SERI 22', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO22', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(40, 230, 'MNCDT20', 'REKSA DANA MNC DANA TERPROTEKSI SERI 20', 'DANA TERPROTEKSI SERI 20', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO20', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(41, 247, 'MNCDT27', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 27', 'DANA TERPROTEKSI SERI 27', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO27', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(42, 248, 'MNCDT26', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 26', 'DANA TERPROTEKSI SERI 26', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO26', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(43, 10248, 'MNCDT28', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 28', 'DANA TERPROTEKSI SERI 28', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO28', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(44, 10249, 'MNCDT24', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 24', 'DANA TERPROTEKSI SERI 24', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO24', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(45, 10250, 'MNCDT25', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 25', 'DANA TERPROTEKSI SERI 25', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO25', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(46, 10251, 'MNCDT30', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 30', 'DANA TERPROTEKSI SERI 30', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO30', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(47, 10252, 'MNC36ETF', 'REKSA DANA ETF MNC36 LIKUID', 'MNC36 LIKUID', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002ETCM36LIK00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(48, 10253, 'MNCDT29', 'REKSA DANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 29', 'DANA TERPROTEKSI SERI 29', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO29', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(49, 10254, 'MNCDT31', 'REKSADANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 31', 'DANA TERPROTEKSI SERI 31', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO31', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(50, 10255, 'MNCDSPT', 'REKSA DANA SYARIAH MNC SYARIAH PENDAPATAN TETAP', 'SYARIAH PENDAPATAN TETAP', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002FISDANTTP00', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(51, 10256, 'MNCDT32', 'REKSADANA TERPROTEKSI MNC DANA TERPROTEKSI SERI 32', 'DANA TERPROTEKSI SERI 32', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.00', '0.0000', '0.0000', NULL, '0.0000', '0.00', NULL, NULL, '0.00', NULL, 'EP002PFCDANPRO32', 0, 0, 0, 0, 0, NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_product_bank_account`
--

CREATE TABLE `ms_product_bank_account` (
  `prod_bankacc_key` int(11) NOT NULL,
  `product_key` int(11) DEFAULT NULL,
  `bank_account_key` int(11) DEFAULT NULL,
  `bank_account_name` varchar(100) NOT NULL,
  `bank_account_purpose` int(11) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_product_category`
--

CREATE TABLE `ms_product_category` (
  `product_category_key` int(11) NOT NULL,
  `category_code` varchar(20) DEFAULT NULL,
  `category_name` varchar(50) DEFAULT NULL,
  `category_desc` varchar(150) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_product_category`
--

INSERT INTO `ms_product_category` (`product_category_key`, `category_code`, `category_name`, `category_desc`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'MF', 'Mutual Fund', 'Reksadana', 10, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'KPD', 'KPD_Discrey', 'KPD_Discrey', 20, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 'DIRE', 'Realestate Fund', 'Realestate', 30, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 'DINFRA', 'Infrastructure Fund', 'Infrastructure', 40, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_product_type`
--

CREATE TABLE `ms_product_type` (
  `product_type_key` int(11) NOT NULL,
  `product_type_code` varchar(20) DEFAULT NULL,
  `product_type_name` varchar(50) DEFAULT NULL,
  `product_type_desc` varchar(150) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_product_type`
--

INSERT INTO `ms_product_type` (`product_type_key`, `product_type_code`, `product_type_name`, `product_type_desc`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'PUBLIC', 'Public Fund', 'Reksa dana yang dijual untuk public', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'PRIVATE', 'Private Fund', 'Reksadana terbatas atau tidak dijual umum', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_risk_profile`
--

CREATE TABLE `ms_risk_profile` (
  `risk_profile_key` int(11) NOT NULL,
  `risk_code` varchar(30) NOT NULL,
  `risk_name` varchar(50) DEFAULT NULL,
  `risk_desc` varchar(1000) DEFAULT NULL,
  `min_score` decimal(8,2) NOT NULL,
  `max_score` decimal(8,2) NOT NULL,
  `max_flag` tinyint(1) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_risk_profile`
--

INSERT INTO `ms_risk_profile` (`risk_profile_key`, `risk_code`, `risk_name`, `risk_desc`, `min_score`, `max_score`, `max_flag`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'Rendah', 'Konservatif', 'Investor mengutamakan tingkat keutuhan nilai pokok investasi dengan risiko fluktuasi investasi yang relatif rendah untuk memenuhi kebutuhan aliran kas bulanan.', '0.00', '0.00', 0, 1, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '1', NULL, NULL),
(2, 'RendahMenengah', 'Konservatif Moderate', 'Investor masih tetap mengutamakan tingkat keutuhan nilai pokok investasi, namun bersedia menerima fluktuasi investasi jangka pendek untuk mendapatkan hasil investasi.', '0.00', '0.00', 0, 2, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2', NULL, NULL),
(3, 'MenengahTinggi', 'Moderate Agresif', 'Investor mulai mencoba alternatif untuk mendapatkan hasil investasi yang relatif tinggi dan risiko fluktuasi investasi yang juga relatif tinggi.', '0.00', '0.00', 0, 3, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '3', NULL, NULL),
(4, 'Tinggi', 'Agresif', 'Investor mengutamakan hasil investasi yang tinggi dalam jangka panjang dan siap menerima tingginya risiko investasi.', '0.00', '0.00', 1, 4, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '4', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ms_securities`
--

CREATE TABLE `ms_securities` (
  `sec_key` int(11) NOT NULL,
  `sec_code` varchar(20) NOT NULL,
  `sec_name` varchar(150) NOT NULL,
  `securities_category` int(11) NOT NULL,
  `security_type` int(11) NOT NULL,
  `date_issued` datetime(6) DEFAULT NULL,
  `date_matured` date DEFAULT NULL,
  `currency_key` int(11) DEFAULT NULL,
  `security_status` int(11) NOT NULL,
  `isin_code` varchar(30) DEFAULT NULL,
  `sec_classification` int(11) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ms_securities_sector`
--

CREATE TABLE `ms_securities_sector` (
  `sector_key` int(11) NOT NULL,
  `sector_code` varchar(20) DEFAULT NULL,
  `sector_name` varchar(150) NOT NULL,
  `sector_description` varchar(255) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `ms_securities_sector`
--

INSERT INTO `ms_securities_sector` (`sector_key`, `sector_code`, `sector_name`, `sector_description`, `rec_order`, `rec_status`, `rec_created_date`, `rec_created_by`, `rec_modified_date`, `rec_modified_by`, `rec_image1`, `rec_image2`, `rec_approval_status`, `rec_approval_stage`, `rec_approved_date`, `rec_approved_by`, `rec_deleted_date`, `rec_deleted_by`, `rec_attribute_id1`, `rec_attribute_id2`, `rec_attribute_id3`) VALUES
(1, 'AGRI', 'AGRI', 'AGRI', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'BASIC-IND', 'BASIC-IND', 'BASIC-IND', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 'CONSUMER', 'CONSUMER', 'CONSUMER', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 'FINANCE', 'FINANCE', 'FINANCE', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 'GOVT', 'GOVT', 'GOVT', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 'INFRA', 'INFRA', 'INFRA', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(7, 'MINING', 'MINING', 'MINING', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(8, 'MISC-IND', 'MISC-IND', 'MISC-IND', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(9, 'PROPERTY', 'PROPERTY', 'PROPERTY', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(10, 'TRADE', 'TRADE', 'TRADE', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(11, 'OTHER', 'Other', 'Other', 0, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `sc_user_login`
--

CREATE TABLE `sc_user_login` (
  `ulogin_key` int(11) NOT NULL,
  `ulogin_name` varchar(30) NOT NULL,
  `ulogin_full_name` varchar(50) NOT NULL,
  `ulogin_password` varchar(30) NOT NULL,
  `ulogin_email` varchar(50) NOT NULL,
  `ulogin_pin` varchar(6) DEFAULT NULL,
  `ulogin_mobileno` varchar(20) DEFAULT NULL,
  `ulogin_role` varchar(50) DEFAULT NULL,
  `ulogin_locked` tinyint(1) NOT NULL,
  `ulogin_enabled` tinyint(1) NOT NULL,
  `ulogin_failed_count` int(11) NOT NULL,
  `ulogin_last_access` datetime(6) NOT NULL,
  `ulogin_email_verified` datetime(6) NOT NULL,
  `ulogin_mobileno_verified` datetime(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_account`
--

CREATE TABLE `tr_account` (
  `acc_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `customer_key` int(11) NOT NULL,
  `account_name` varchar(100) NOT NULL,
  `account_no` varchar(30) NOT NULL,
  `ifua_no` varchar(20) DEFAULT NULL,
  `ifua_name` varchar(100) DEFAULT NULL,
  `acc_status` int(11) NOT NULL,
  `acc_suspend_flag` tinyint(1) NOT NULL,
  `acc_suspend_modified_date` datetime(6) NOT NULL,
  `acc_suspend_reason` varchar(100) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_account_agent`
--

CREATE TABLE `tr_account_agent` (
  `aca_key` int(11) NOT NULL,
  `acc_key` int(11) NOT NULL,
  `eff_date` datetime(6) NOT NULL,
  `branch_key` int(11) DEFAULT NULL,
  `agent_key` int(11) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_balance`
--

CREATE TABLE `tr_balance` (
  `balance_key` int(11) NOT NULL,
  `ag_key` int(11) NOT NULL,
  `transaction_key` int(11) NOT NULL,
  `trans_date` date DEFAULT NULL,
  `trx_code` int(11) NOT NULL,
  `available_unit` decimal(12,4) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_currency_rate`
--

CREATE TABLE `tr_currency_rate` (
  `curr_rate_key` int(11) NOT NULL,
  `rate_date` date NOT NULL,
  `rate_type` int(11) NOT NULL,
  `rate_value` decimal(9,2) NOT NULL,
  `currency_key` int(11) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_nav`
--

CREATE TABLE `tr_nav` (
  `nav_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `nav_date` date NOT NULL,
  `nav_value` decimal(18,4) NOT NULL,
  `original_value` decimal(18,6) NOT NULL,
  `nav_status` int(11) NOT NULL,
  `prod_aum_total` decimal(18,4) NOT NULL,
  `prod_unit_total` decimal(18,4) NOT NULL,
  `publish_mode` int(11) NOT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_prospective_customer`
--

CREATE TABLE `tr_prospective_customer` (
  `prospective_key` int(11) NOT NULL,
  `agent_key` int(11) NOT NULL,
  `date_appointment` date NOT NULL,
  `pro_name` varchar(50) NOT NULL,
  `pro_email_address` varchar(50) DEFAULT NULL,
  `pro_honeno` varchar(20) NOT NULL,
  `pro_phoneno_alt` varchar(20) DEFAULT NULL,
  `prospect_city` varchar(50) NOT NULL,
  `prospect_notes` varchar(200) DEFAULT NULL,
  `closed_date` date NOT NULL,
  `closed_notes` varchar(200) DEFAULT NULL,
  `cif_number` varchar(20) DEFAULT NULL,
  `prospect_status` int(11) NOT NULL,
  `pro_registration_link` varchar(255) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_transaction`
--

CREATE TABLE `tr_transaction` (
  `transaction_key` int(11) NOT NULL,
  `id_transaction` int(11) DEFAULT NULL,
  `branch_key` int(11) DEFAULT NULL,
  `agent_key` int(11) DEFAULT NULL,
  `customer_key` int(11) NOT NULL,
  `product_key` int(11) NOT NULL,
  `trans_status_key` int(11) NOT NULL,
  `trans_date` datetime(6) NOT NULL,
  `trans_type_key` int(11) NOT NULL,
  `trx_code` varchar(5) NOT NULL,
  `nav_date` date NOT NULL,
  `entry_mode` varchar(20) NOT NULL,
  `trans_amount` decimal(12,2) NOT NULL,
  `trans_unit` decimal(12,6) NOT NULL,
  `trans_unit_percent` decimal(5,2) DEFAULT NULL,
  `flag_redempt_all` tinyint(1) NOT NULL,
  `flag_newsub` tinyint(1) NOT NULL,
  `trans_fee_percent` decimal(5,2) NOT NULL,
  `trans_fee_amount` decimal(9,2) NOT NULL,
  `charges_fee_amount` decimal(9,2) NOT NULL,
  `services_fee_amount` decimal(9,2) NOT NULL,
  `total_amount` decimal(12,2) NOT NULL,
  `settlement_date` date DEFAULT NULL,
  `trans_bank_accno` varchar(20) DEFAULT NULL,
  `trans_bankacc_name` varchar(50) DEFAULT NULL,
  `trans_bank_key` int(11) DEFAULT NULL,
  `trans_remarks` varchar(150) DEFAULT NULL,
  `trans_references` varchar(150) DEFAULT NULL,
  `promo_code` varchar(20) DEFAULT NULL,
  `risk_waiver` tinyint(1) NOT NULL,
  `addto_auto_invest` tinyint(1) NOT NULL,
  `trans_source` int(11) NOT NULL,
  `proceed_date` datetime DEFAULT NULL,
  `proceed_amount` decimal(12,4) NOT NULL,
  `sent_date` datetime DEFAULT NULL,
  `sent_references` varchar(50) DEFAULT NULL,
  `confirmed_date` datetime DEFAULT NULL,
  `posted_date` datetime DEFAULT NULL,
  `posted_units` decimal(12,6) NOT NULL,
  `aca_key` int(11) DEFAULT NULL,
  `settled_date` date DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_transaction_confirmation`
--

CREATE TABLE `tr_transaction_confirmation` (
  `tc_key` int(11) NOT NULL,
  `confirm_date` date NOT NULL,
  `transaction_key` int(11) NOT NULL,
  `confirmed_amount` decimal(12,4) NOT NULL,
  `confirmed_unit` decimal(12,4) NOT NULL,
  `confirm_result` int(11) NOT NULL,
  `confirmed_amount_diff` decimal(12,4) NOT NULL,
  `confirmed_unit_diff` decimal(12,4) NOT NULL,
  `confirmed_remarks` varchar(100) DEFAULT NULL,
  `confirmed_references` varchar(100) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_transaction_fifo`
--

CREATE TABLE `tr_transaction_fifo` (
  `trans_fifo_key` int(11) NOT NULL,
  `trans_red_key` int(11) DEFAULT NULL,
  `trans_sub_key` int(11) DEFAULT NULL,
  `sub_aca_key` int(11) DEFAULT NULL,
  `holding_days` int(11) DEFAULT NULL,
  `trans_unit` decimal(12,4) NOT NULL,
  `fee_nav_mode` int(11) NOT NULL,
  `trans_amount` decimal(12,4) NOT NULL,
  `trans_fee_amount` decimal(12,4) NOT NULL,
  `trans_fee_tax` decimal(12,4) NOT NULL,
  `trans_nett_amount` decimal(12,4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_transaction_settlement`
--

CREATE TABLE `tr_transaction_settlement` (
  `settlement_key` int(11) NOT NULL,
  `transaction_key` int(11) DEFAULT NULL,
  `settle_purposed` varchar(50) NOT NULL,
  `settle_date` date NOT NULL,
  `settle_nominal` decimal(12,4) NOT NULL,
  `settled_status` int(11) NOT NULL,
  `settle_realized_date` date NOT NULL,
  `settle_remarks` varchar(100) DEFAULT NULL,
  `settle_references` varchar(100) DEFAULT NULL,
  `source_bank_account_key` int(11) DEFAULT NULL,
  `target_bank_account_key` int(11) NOT NULL,
  `settle_notes` varchar(100) DEFAULT NULL,
  `settle_acknowledgement` varchar(100) DEFAULT NULL,
  `rec_order` int(11) DEFAULT NULL,
  `rec_status` tinyint(3) UNSIGNED NOT NULL,
  `rec_created_date` datetime(6) DEFAULT NULL,
  `rec_created_by` varchar(30) DEFAULT NULL,
  `rec_modified_date` datetime(6) DEFAULT NULL,
  `rec_modified_by` varchar(30) DEFAULT NULL,
  `rec_image1` varchar(50) DEFAULT NULL,
  `rec_image2` varchar(50) DEFAULT NULL,
  `rec_approval_status` tinyint(3) UNSIGNED DEFAULT NULL,
  `rec_approval_stage` int(11) DEFAULT NULL,
  `rec_approved_date` datetime(6) DEFAULT NULL,
  `rec_approved_by` varchar(30) DEFAULT NULL,
  `rec_deleted_date` datetime(6) DEFAULT NULL,
  `rec_deleted_by` varchar(30) DEFAULT NULL,
  `rec_attribute_id1` varchar(30) DEFAULT NULL,
  `rec_attribute_id2` varchar(30) DEFAULT NULL,
  `rec_attribute_id3` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `tr_transaction_status`
--

CREATE TABLE `tr_transaction_status` (
  `trans_status_key` int(11) NOT NULL,
  `status_code` longtext DEFAULT NULL,
  `status_description` longtext DEFAULT NULL,
  `status_order` int(11) NOT NULL,
  `status_phase` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `tr_transaction_status`
--

INSERT INTO `tr_transaction_status` (`trans_status_key`, `status_code`, `status_description`, `status_order`, `status_phase`) VALUES
(1, 'ENTRIED', 'Transaction has been successfully entry and still in progress', 10, 'STEP-1'),
(2, 'SUBMITTED', 'transaction has been approved, means it has been submitted to Back Office', 11, 'STEP-1'),
(3, 'CORRECTED', 'Transaction is still open for correction or update dan in still progress ', 12, 'STEP-1'),
(4, 'DELETED', 'The transaction has been deleted', 13, 'STEP-1'),
(5, 'PROCEED', 'Transaction has been cutoff proceed. After Submitted.', 20, 'STEP-2'),
(6, 'BATCHED', 'Transaction has been batched or sent to S-INVEST', 30, 'STEP-2'),
(7, 'CONFIRMED', 'Transaction has been confirmed by Custodian', 40, 'STEP-3'),
(8, 'POSTED', 'Transaction has been posted to update the holding unit (balance)', 50, 'STEP-3'),
(9, 'SETTLED', 'Transaction has been settled (when need payment process)', 60, 'STEP-3');

-- --------------------------------------------------------

--
-- Table structure for table `tr_transaction_type`
--

CREATE TABLE `tr_transaction_type` (
  `trans_type_key` int(11) NOT NULL,
  `type_code` longtext DEFAULT NULL,
  `type_description` longtext DEFAULT NULL,
  `type_order` int(11) NOT NULL,
  `type_domain` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `tr_transaction_type`
--

INSERT INTO `tr_transaction_type` (`trans_type_key`, `type_code`, `type_description`, `type_order`, `type_domain`) VALUES
(1, 'SUB', 'Subscription', 1, 'FRONT'),
(2, 'RED', 'Redeemption', 6, 'FRONT'),
(3, 'SWTOT', 'Switch Out', 5, 'FRONT'),
(4, 'SWTIN', 'Switch In', 2, 'FRONT'),
(5, 'DIVUNIT', 'Dividen Unit', 3, 'BACK'),
(6, 'DIVCASH', 'Dividen Cash', 4, 'BACK'),
(7, 'MASSRED', 'Mass Redemption', 7, 'BACK'),
(8, 'LIQUIDATE', 'Maturity', 8, 'BACK');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `Id` int(11) NOT NULL,
  `Username` longtext DEFAULT NULL,
  `Password` longtext DEFAULT NULL,
  `Role` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `__efmigrationshistory`
--

CREATE TABLE `__efmigrationshistory` (
  `MigrationId` varchar(95) NOT NULL,
  `ProductVersion` varchar(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `__efmigrationshistory`
--

INSERT INTO `__efmigrationshistory` (`MigrationId`, `ProductVersion`) VALUES
('20200722132939_MamCore_InitDB', '3.1.6');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `ffs_alloc_instrument`
--
ALTER TABLE `ffs_alloc_instrument`
  ADD PRIMARY KEY (`alloc_instrument_key`),
  ADD KEY `IX_ffs_alloc_instrument_instrument_key` (`instrument_key`),
  ADD KEY `IX_ffs_alloc_instrument_periode_key` (`periode_key`),
  ADD KEY `IX_ffs_alloc_instrument_product_key` (`product_key`);

--
-- Indexes for table `ffs_alloc_sector`
--
ALTER TABLE `ffs_alloc_sector`
  ADD PRIMARY KEY (`alloc_sector_key`),
  ADD KEY `IX_ffs_alloc_sector_periode_key` (`periode_key`),
  ADD KEY `IX_ffs_alloc_sector_product_key` (`product_key`),
  ADD KEY `IX_ffs_alloc_sector_sector_key` (`sector_key`);

--
-- Indexes for table `ffs_alloc_security`
--
ALTER TABLE `ffs_alloc_security`
  ADD PRIMARY KEY (`alloc_security_key`),
  ADD KEY `IX_ffs_alloc_security_periode_key` (`periode_key`),
  ADD KEY `IX_ffs_alloc_security_product_key` (`product_key`),
  ADD KEY `IX_ffs_alloc_security_sec_key` (`sec_key`);

--
-- Indexes for table `ffs_nav_performance`
--
ALTER TABLE `ffs_nav_performance`
  ADD PRIMARY KEY (`nav_perform_key`),
  ADD KEY `IX_ffs_nav_performance_periode_key` (`periode_key`),
  ADD KEY `IX_ffs_nav_performance_product_key` (`product_key`);

--
-- Indexes for table `ffs_periode`
--
ALTER TABLE `ffs_periode`
  ADD PRIMARY KEY (`periode_key`);

--
-- Indexes for table `ffs_publish`
--
ALTER TABLE `ffs_publish`
  ADD PRIMARY KEY (`ffs_key`),
  ADD KEY `IX_ffs_publish_periode_key` (`periode_key`),
  ADD KEY `IX_ffs_publish_product_key` (`product_key`);

--
-- Indexes for table `ms_agent`
--
ALTER TABLE `ms_agent`
  ADD PRIMARY KEY (`agent_key`);

--
-- Indexes for table `ms_agent_agreement`
--
ALTER TABLE `ms_agent_agreement`
  ADD PRIMARY KEY (`agreement_key`),
  ADD KEY `IX_ms_agent_agreement_branch_key` (`branch_key`);

--
-- Indexes for table `ms_agent_bank_account`
--
ALTER TABLE `ms_agent_bank_account`
  ADD PRIMARY KEY (`agent_bankacc_key`),
  ADD KEY `IX_ms_agent_bank_account_agent_key` (`agent_key`),
  ADD KEY `IX_ms_agent_bank_account_bank_account_key` (`bank_account_key`);

--
-- Indexes for table `ms_agent_branch`
--
ALTER TABLE `ms_agent_branch`
  ADD PRIMARY KEY (`agent_branch_key`),
  ADD KEY `IX_ms_agent_branch_agent_key` (`agent_key`),
  ADD KEY `IX_ms_agent_branch_branch_key` (`branch_key`);

--
-- Indexes for table `ms_agent_customer`
--
ALTER TABLE `ms_agent_customer`
  ADD PRIMARY KEY (`agent_customer_key`),
  ADD KEY `IX_ms_agent_customer_agent_key` (`agent_key`),
  ADD KEY `IX_ms_agent_customer_customer_key` (`customer_key`);

--
-- Indexes for table `ms_agent_detail`
--
ALTER TABLE `ms_agent_detail`
  ADD PRIMARY KEY (`agent_key`),
  ADD KEY `IX_ms_agent_detail_country_key` (`country_key`);

--
-- Indexes for table `ms_agent_license`
--
ALTER TABLE `ms_agent_license`
  ADD PRIMARY KEY (`aglic_key`),
  ADD KEY `IX_ms_agent_license_agent_key` (`agent_key`);

--
-- Indexes for table `ms_agent_product`
--
ALTER TABLE `ms_agent_product`
  ADD PRIMARY KEY (`agent_product_key`),
  ADD KEY `IX_ms_agent_product_branch_key` (`branch_key`),
  ADD KEY `IX_ms_agent_product_product_key` (`product_key`);

--
-- Indexes for table `ms_bank`
--
ALTER TABLE `ms_bank`
  ADD PRIMARY KEY (`bank_key`);

--
-- Indexes for table `ms_bank_account`
--
ALTER TABLE `ms_bank_account`
  ADD PRIMARY KEY (`bank_account_key`),
  ADD KEY `IX_ms_bank_account_bank_key` (`bank_key`),
  ADD KEY `IX_ms_bank_account_currency_key` (`currency_key`);

--
-- Indexes for table `ms_branch`
--
ALTER TABLE `ms_branch`
  ADD PRIMARY KEY (`branch_key`),
  ADD KEY `IX_ms_branch_city_key` (`city_key`),
  ADD KEY `IX_ms_branch_participant_key` (`participant_key`);

--
-- Indexes for table `ms_city`
--
ALTER TABLE `ms_city`
  ADD PRIMARY KEY (`city_key`),
  ADD KEY `IX_ms_city_parent_key` (`parent_key`);

--
-- Indexes for table `ms_country`
--
ALTER TABLE `ms_country`
  ADD PRIMARY KEY (`country_key`),
  ADD KEY `IX_ms_country_currency_key` (`currency_key`);

--
-- Indexes for table `ms_currency`
--
ALTER TABLE `ms_currency`
  ADD PRIMARY KEY (`currency_key`);

--
-- Indexes for table `ms_custodian_bank`
--
ALTER TABLE `ms_custodian_bank`
  ADD PRIMARY KEY (`custodian_key`);

--
-- Indexes for table `ms_customer`
--
ALTER TABLE `ms_customer`
  ADD PRIMARY KEY (`customer_key`),
  ADD KEY `IX_ms_customer_closeacc_agent_key` (`closeacc_agent_key`),
  ADD KEY `IX_ms_customer_closeacc_branch_key` (`closeacc_branch_key`),
  ADD KEY `IX_ms_customer_openacc_agent_key` (`openacc_agent_key`),
  ADD KEY `IX_ms_customer_openacc_branch_key` (`openacc_branch_key`),
  ADD KEY `IX_ms_customer_participant_key` (`participant_key`);

--
-- Indexes for table `ms_customer_login`
--
ALTER TABLE `ms_customer_login`
  ADD PRIMARY KEY (`cust_login_key`),
  ADD KEY `IX_ms_customer_login_customer_key` (`customer_key`);

--
-- Indexes for table `ms_cutomer_detail`
--
ALTER TABLE `ms_cutomer_detail`
  ADD PRIMARY KEY (`customer_key`);

--
-- Indexes for table `ms_file`
--
ALTER TABLE `ms_file`
  ADD PRIMARY KEY (`file_key`);

--
-- Indexes for table `ms_fund_structure`
--
ALTER TABLE `ms_fund_structure`
  ADD PRIMARY KEY (`fund_structure_key`);

--
-- Indexes for table `ms_fund_type`
--
ALTER TABLE `ms_fund_type`
  ADD PRIMARY KEY (`fund_type_key`);

--
-- Indexes for table `ms_geolocation`
--
ALTER TABLE `ms_geolocation`
  ADD PRIMARY KEY (`location_key`);

--
-- Indexes for table `ms_holiday`
--
ALTER TABLE `ms_holiday`
  ADD PRIMARY KEY (`holiday_key`);

--
-- Indexes for table `ms_instrument`
--
ALTER TABLE `ms_instrument`
  ADD PRIMARY KEY (`instrument_key`);

--
-- Indexes for table `ms_participant`
--
ALTER TABLE `ms_participant`
  ADD PRIMARY KEY (`participant_key`);

--
-- Indexes for table `ms_product`
--
ALTER TABLE `ms_product`
  ADD PRIMARY KEY (`product_key`),
  ADD KEY `IX_ms_product_currency_key` (`currency_key`),
  ADD KEY `IX_ms_product_custodian_key` (`custodian_key`),
  ADD KEY `IX_ms_product_fund_structure_key` (`fund_structure_key`),
  ADD KEY `IX_ms_product_fund_type_key` (`fund_type_key`),
  ADD KEY `IX_ms_product_product_category_key` (`product_category_key`),
  ADD KEY `IX_ms_product_product_type_key` (`product_type_key`),
  ADD KEY `IX_ms_product_risk_profile_key` (`risk_profile_key`);

--
-- Indexes for table `ms_product_bank_account`
--
ALTER TABLE `ms_product_bank_account`
  ADD PRIMARY KEY (`prod_bankacc_key`),
  ADD KEY `IX_ms_product_bank_account_bank_account_key` (`bank_account_key`),
  ADD KEY `IX_ms_product_bank_account_product_key` (`product_key`);

--
-- Indexes for table `ms_product_category`
--
ALTER TABLE `ms_product_category`
  ADD PRIMARY KEY (`product_category_key`);

--
-- Indexes for table `ms_product_type`
--
ALTER TABLE `ms_product_type`
  ADD PRIMARY KEY (`product_type_key`);

--
-- Indexes for table `ms_risk_profile`
--
ALTER TABLE `ms_risk_profile`
  ADD PRIMARY KEY (`risk_profile_key`);

--
-- Indexes for table `ms_securities`
--
ALTER TABLE `ms_securities`
  ADD PRIMARY KEY (`sec_key`),
  ADD KEY `IX_ms_securities_currency_key` (`currency_key`);

--
-- Indexes for table `ms_securities_sector`
--
ALTER TABLE `ms_securities_sector`
  ADD PRIMARY KEY (`sector_key`);

--
-- Indexes for table `sc_user_login`
--
ALTER TABLE `sc_user_login`
  ADD PRIMARY KEY (`ulogin_key`);

--
-- Indexes for table `tr_account`
--
ALTER TABLE `tr_account`
  ADD PRIMARY KEY (`acc_key`),
  ADD KEY `IX_tr_account_customer_key` (`customer_key`),
  ADD KEY `IX_tr_account_product_key` (`product_key`);

--
-- Indexes for table `tr_account_agent`
--
ALTER TABLE `tr_account_agent`
  ADD PRIMARY KEY (`aca_key`),
  ADD KEY `IX_tr_account_agent_acc_key` (`acc_key`),
  ADD KEY `IX_tr_account_agent_agent_key` (`agent_key`),
  ADD KEY `IX_tr_account_agent_branch_key` (`branch_key`);

--
-- Indexes for table `tr_balance`
--
ALTER TABLE `tr_balance`
  ADD PRIMARY KEY (`balance_key`),
  ADD KEY `IX_tr_balance_ag_key` (`ag_key`);

--
-- Indexes for table `tr_currency_rate`
--
ALTER TABLE `tr_currency_rate`
  ADD PRIMARY KEY (`curr_rate_key`),
  ADD KEY `IX_tr_currency_rate_currency_key` (`currency_key`);

--
-- Indexes for table `tr_nav`
--
ALTER TABLE `tr_nav`
  ADD PRIMARY KEY (`nav_key`),
  ADD KEY `IX_tr_nav_product_key` (`product_key`);

--
-- Indexes for table `tr_prospective_customer`
--
ALTER TABLE `tr_prospective_customer`
  ADD PRIMARY KEY (`prospective_key`),
  ADD KEY `IX_tr_prospective_customer_agent_key` (`agent_key`);

--
-- Indexes for table `tr_transaction`
--
ALTER TABLE `tr_transaction`
  ADD PRIMARY KEY (`transaction_key`),
  ADD KEY `IX_tr_transaction_aca_key` (`aca_key`),
  ADD KEY `IX_tr_transaction_agent_key` (`agent_key`),
  ADD KEY `IX_tr_transaction_branch_key` (`branch_key`),
  ADD KEY `IX_tr_transaction_customer_key` (`customer_key`),
  ADD KEY `IX_tr_transaction_product_key` (`product_key`),
  ADD KEY `IX_tr_transaction_trans_bank_key` (`trans_bank_key`),
  ADD KEY `IX_tr_transaction_trans_status_key` (`trans_status_key`),
  ADD KEY `IX_tr_transaction_trans_type_key` (`trans_type_key`);

--
-- Indexes for table `tr_transaction_confirmation`
--
ALTER TABLE `tr_transaction_confirmation`
  ADD PRIMARY KEY (`tc_key`),
  ADD UNIQUE KEY `IX_tr_transaction_confirmation_transaction_key` (`transaction_key`);

--
-- Indexes for table `tr_transaction_fifo`
--
ALTER TABLE `tr_transaction_fifo`
  ADD PRIMARY KEY (`trans_fifo_key`),
  ADD KEY `IX_tr_transaction_fifo_sub_aca_key` (`sub_aca_key`),
  ADD KEY `IX_tr_transaction_fifo_trans_red_key` (`trans_red_key`),
  ADD KEY `IX_tr_transaction_fifo_trans_sub_key` (`trans_sub_key`);

--
-- Indexes for table `tr_transaction_settlement`
--
ALTER TABLE `tr_transaction_settlement`
  ADD PRIMARY KEY (`settlement_key`),
  ADD KEY `IX_tr_transaction_settlement_source_bank_account_key` (`source_bank_account_key`),
  ADD KEY `IX_tr_transaction_settlement_target_bank_account_key` (`target_bank_account_key`),
  ADD KEY `IX_tr_transaction_settlement_transaction_key` (`transaction_key`);

--
-- Indexes for table `tr_transaction_status`
--
ALTER TABLE `tr_transaction_status`
  ADD PRIMARY KEY (`trans_status_key`);

--
-- Indexes for table `tr_transaction_type`
--
ALTER TABLE `tr_transaction_type`
  ADD PRIMARY KEY (`trans_type_key`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `__efmigrationshistory`
--
ALTER TABLE `__efmigrationshistory`
  ADD PRIMARY KEY (`MigrationId`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `ffs_alloc_instrument`
--
ALTER TABLE `ffs_alloc_instrument`
  MODIFY `alloc_instrument_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ffs_alloc_sector`
--
ALTER TABLE `ffs_alloc_sector`
  MODIFY `alloc_sector_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ffs_alloc_security`
--
ALTER TABLE `ffs_alloc_security`
  MODIFY `alloc_security_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ffs_nav_performance`
--
ALTER TABLE `ffs_nav_performance`
  MODIFY `nav_perform_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ffs_periode`
--
ALTER TABLE `ffs_periode`
  MODIFY `periode_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ffs_publish`
--
ALTER TABLE `ffs_publish`
  MODIFY `ffs_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_agent`
--
ALTER TABLE `ms_agent`
  MODIFY `agent_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=41;

--
-- AUTO_INCREMENT for table `ms_agent_agreement`
--
ALTER TABLE `ms_agent_agreement`
  MODIFY `agreement_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_agent_bank_account`
--
ALTER TABLE `ms_agent_bank_account`
  MODIFY `agent_bankacc_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_agent_branch`
--
ALTER TABLE `ms_agent_branch`
  MODIFY `agent_branch_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `ms_agent_customer`
--
ALTER TABLE `ms_agent_customer`
  MODIFY `agent_customer_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `ms_agent_detail`
--
ALTER TABLE `ms_agent_detail`
  MODIFY `agent_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_agent_license`
--
ALTER TABLE `ms_agent_license`
  MODIFY `aglic_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_agent_product`
--
ALTER TABLE `ms_agent_product`
  MODIFY `agent_product_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_bank`
--
ALTER TABLE `ms_bank`
  MODIFY `bank_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_bank_account`
--
ALTER TABLE `ms_bank_account`
  MODIFY `bank_account_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_branch`
--
ALTER TABLE `ms_branch`
  MODIFY `branch_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=20;

--
-- AUTO_INCREMENT for table `ms_city`
--
ALTER TABLE `ms_city`
  MODIFY `city_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_country`
--
ALTER TABLE `ms_country`
  MODIFY `country_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=242;

--
-- AUTO_INCREMENT for table `ms_currency`
--
ALTER TABLE `ms_currency`
  MODIFY `currency_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `ms_custodian_bank`
--
ALTER TABLE `ms_custodian_bank`
  MODIFY `custodian_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `ms_customer`
--
ALTER TABLE `ms_customer`
  MODIFY `customer_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `ms_customer_login`
--
ALTER TABLE `ms_customer_login`
  MODIFY `cust_login_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_cutomer_detail`
--
ALTER TABLE `ms_cutomer_detail`
  MODIFY `customer_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_file`
--
ALTER TABLE `ms_file`
  MODIFY `file_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_fund_structure`
--
ALTER TABLE `ms_fund_structure`
  MODIFY `fund_structure_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `ms_fund_type`
--
ALTER TABLE `ms_fund_type`
  MODIFY `fund_type_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `ms_geolocation`
--
ALTER TABLE `ms_geolocation`
  MODIFY `location_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_holiday`
--
ALTER TABLE `ms_holiday`
  MODIFY `holiday_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_instrument`
--
ALTER TABLE `ms_instrument`
  MODIFY `instrument_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `ms_participant`
--
ALTER TABLE `ms_participant`
  MODIFY `participant_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=401;

--
-- AUTO_INCREMENT for table `ms_product`
--
ALTER TABLE `ms_product`
  MODIFY `product_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=52;

--
-- AUTO_INCREMENT for table `ms_product_bank_account`
--
ALTER TABLE `ms_product_bank_account`
  MODIFY `prod_bankacc_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_product_category`
--
ALTER TABLE `ms_product_category`
  MODIFY `product_category_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `ms_product_type`
--
ALTER TABLE `ms_product_type`
  MODIFY `product_type_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `ms_risk_profile`
--
ALTER TABLE `ms_risk_profile`
  MODIFY `risk_profile_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `ms_securities`
--
ALTER TABLE `ms_securities`
  MODIFY `sec_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ms_securities_sector`
--
ALTER TABLE `ms_securities_sector`
  MODIFY `sector_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- AUTO_INCREMENT for table `sc_user_login`
--
ALTER TABLE `sc_user_login`
  MODIFY `ulogin_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_account`
--
ALTER TABLE `tr_account`
  MODIFY `acc_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_account_agent`
--
ALTER TABLE `tr_account_agent`
  MODIFY `aca_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_balance`
--
ALTER TABLE `tr_balance`
  MODIFY `balance_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_currency_rate`
--
ALTER TABLE `tr_currency_rate`
  MODIFY `curr_rate_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_nav`
--
ALTER TABLE `tr_nav`
  MODIFY `nav_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_prospective_customer`
--
ALTER TABLE `tr_prospective_customer`
  MODIFY `prospective_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_transaction`
--
ALTER TABLE `tr_transaction`
  MODIFY `transaction_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_transaction_confirmation`
--
ALTER TABLE `tr_transaction_confirmation`
  MODIFY `tc_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_transaction_fifo`
--
ALTER TABLE `tr_transaction_fifo`
  MODIFY `trans_fifo_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_transaction_settlement`
--
ALTER TABLE `tr_transaction_settlement`
  MODIFY `settlement_key` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tr_transaction_status`
--
ALTER TABLE `tr_transaction_status`
  MODIFY `trans_status_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- AUTO_INCREMENT for table `tr_transaction_type`
--
ALTER TABLE `tr_transaction_type`
  MODIFY `trans_type_key` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `ffs_alloc_instrument`
--
ALTER TABLE `ffs_alloc_instrument`
  ADD CONSTRAINT `FK_ffs_alloc_instrument_ffs_periode_periode_key` FOREIGN KEY (`periode_key`) REFERENCES `ffs_periode` (`periode_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_alloc_instrument_ms_instrument_instrument_key` FOREIGN KEY (`instrument_key`) REFERENCES `ms_instrument` (`instrument_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_alloc_instrument_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE;

--
-- Constraints for table `ffs_alloc_sector`
--
ALTER TABLE `ffs_alloc_sector`
  ADD CONSTRAINT `FK_ffs_alloc_sector_ffs_periode_periode_key` FOREIGN KEY (`periode_key`) REFERENCES `ffs_periode` (`periode_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_alloc_sector_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_alloc_sector_ms_securities_sector_sector_key` FOREIGN KEY (`sector_key`) REFERENCES `ms_securities_sector` (`sector_key`) ON DELETE CASCADE;

--
-- Constraints for table `ffs_alloc_security`
--
ALTER TABLE `ffs_alloc_security`
  ADD CONSTRAINT `FK_ffs_alloc_security_ffs_periode_periode_key` FOREIGN KEY (`periode_key`) REFERENCES `ffs_periode` (`periode_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_alloc_security_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_alloc_security_ms_securities_sec_key` FOREIGN KEY (`sec_key`) REFERENCES `ms_securities` (`sec_key`) ON DELETE CASCADE;

--
-- Constraints for table `ffs_nav_performance`
--
ALTER TABLE `ffs_nav_performance`
  ADD CONSTRAINT `FK_ffs_nav_performance_ffs_periode_periode_key` FOREIGN KEY (`periode_key`) REFERENCES `ffs_periode` (`periode_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_nav_performance_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE;

--
-- Constraints for table `ffs_publish`
--
ALTER TABLE `ffs_publish`
  ADD CONSTRAINT `FK_ffs_publish_ffs_periode_periode_key` FOREIGN KEY (`periode_key`) REFERENCES `ffs_periode` (`periode_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ffs_publish_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_agent_agreement`
--
ALTER TABLE `ms_agent_agreement`
  ADD CONSTRAINT `FK_ms_agent_agreement_ms_branch_branch_key` FOREIGN KEY (`branch_key`) REFERENCES `ms_branch` (`branch_key`);

--
-- Constraints for table `ms_agent_bank_account`
--
ALTER TABLE `ms_agent_bank_account`
  ADD CONSTRAINT `FK_ms_agent_bank_account_ms_agent_agent_key` FOREIGN KEY (`agent_key`) REFERENCES `ms_agent` (`agent_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ms_agent_bank_account_ms_bank_account_bank_account_key` FOREIGN KEY (`bank_account_key`) REFERENCES `ms_bank_account` (`bank_account_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_agent_branch`
--
ALTER TABLE `ms_agent_branch`
  ADD CONSTRAINT `FK_ms_agent_branch_ms_agent_agent_key` FOREIGN KEY (`agent_key`) REFERENCES `ms_agent` (`agent_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ms_agent_branch_ms_branch_branch_key` FOREIGN KEY (`branch_key`) REFERENCES `ms_branch` (`branch_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_agent_customer`
--
ALTER TABLE `ms_agent_customer`
  ADD CONSTRAINT `FK_ms_agent_customer_ms_agent_agent_key` FOREIGN KEY (`agent_key`) REFERENCES `ms_agent` (`agent_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ms_agent_customer_ms_customer_customer_key` FOREIGN KEY (`customer_key`) REFERENCES `ms_customer` (`customer_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_agent_detail`
--
ALTER TABLE `ms_agent_detail`
  ADD CONSTRAINT `FK_ms_agent_detail_ms_country_country_key` FOREIGN KEY (`country_key`) REFERENCES `ms_country` (`country_key`);

--
-- Constraints for table `ms_agent_license`
--
ALTER TABLE `ms_agent_license`
  ADD CONSTRAINT `FK_ms_agent_license_ms_agent_agent_key` FOREIGN KEY (`agent_key`) REFERENCES `ms_agent` (`agent_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_agent_product`
--
ALTER TABLE `ms_agent_product`
  ADD CONSTRAINT `FK_ms_agent_product_ms_branch_branch_key` FOREIGN KEY (`branch_key`) REFERENCES `ms_branch` (`branch_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ms_agent_product_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_bank_account`
--
ALTER TABLE `ms_bank_account`
  ADD CONSTRAINT `FK_ms_bank_account_ms_bank_bank_key` FOREIGN KEY (`bank_key`) REFERENCES `ms_bank` (`bank_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_ms_bank_account_ms_currency_currency_key` FOREIGN KEY (`currency_key`) REFERENCES `ms_currency` (`currency_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_branch`
--
ALTER TABLE `ms_branch`
  ADD CONSTRAINT `FK_ms_branch_ms_city_city_key` FOREIGN KEY (`city_key`) REFERENCES `ms_city` (`city_key`),
  ADD CONSTRAINT `FK_ms_branch_ms_participant_participant_key` FOREIGN KEY (`participant_key`) REFERENCES `ms_participant` (`participant_key`);

--
-- Constraints for table `ms_city`
--
ALTER TABLE `ms_city`
  ADD CONSTRAINT `FK_ms_city_ms_city_parent_key` FOREIGN KEY (`parent_key`) REFERENCES `ms_city` (`city_key`);

--
-- Constraints for table `ms_country`
--
ALTER TABLE `ms_country`
  ADD CONSTRAINT `FK_ms_country_ms_currency_currency_key` FOREIGN KEY (`currency_key`) REFERENCES `ms_currency` (`currency_key`);

--
-- Constraints for table `ms_customer`
--
ALTER TABLE `ms_customer`
  ADD CONSTRAINT `FK_ms_customer_ms_agent_closeacc_agent_key` FOREIGN KEY (`closeacc_agent_key`) REFERENCES `ms_agent` (`agent_key`),
  ADD CONSTRAINT `FK_ms_customer_ms_agent_openacc_agent_key` FOREIGN KEY (`openacc_agent_key`) REFERENCES `ms_agent` (`agent_key`),
  ADD CONSTRAINT `FK_ms_customer_ms_branch_closeacc_branch_key` FOREIGN KEY (`closeacc_branch_key`) REFERENCES `ms_branch` (`branch_key`),
  ADD CONSTRAINT `FK_ms_customer_ms_branch_openacc_branch_key` FOREIGN KEY (`openacc_branch_key`) REFERENCES `ms_branch` (`branch_key`),
  ADD CONSTRAINT `FK_ms_customer_ms_participant_participant_key` FOREIGN KEY (`participant_key`) REFERENCES `ms_participant` (`participant_key`);

--
-- Constraints for table `ms_customer_login`
--
ALTER TABLE `ms_customer_login`
  ADD CONSTRAINT `FK_ms_customer_login_ms_customer_customer_key` FOREIGN KEY (`customer_key`) REFERENCES `ms_customer` (`customer_key`) ON DELETE CASCADE;

--
-- Constraints for table `ms_product`
--
ALTER TABLE `ms_product`
  ADD CONSTRAINT `FK_ms_product_ms_currency_currency_key` FOREIGN KEY (`currency_key`) REFERENCES `ms_currency` (`currency_key`),
  ADD CONSTRAINT `FK_ms_product_ms_custodian_bank_custodian_key` FOREIGN KEY (`custodian_key`) REFERENCES `ms_custodian_bank` (`custodian_key`),
  ADD CONSTRAINT `FK_ms_product_ms_fund_structure_fund_structure_key` FOREIGN KEY (`fund_structure_key`) REFERENCES `ms_fund_structure` (`fund_structure_key`),
  ADD CONSTRAINT `FK_ms_product_ms_fund_type_fund_type_key` FOREIGN KEY (`fund_type_key`) REFERENCES `ms_fund_type` (`fund_type_key`),
  ADD CONSTRAINT `FK_ms_product_ms_product_category_product_category_key` FOREIGN KEY (`product_category_key`) REFERENCES `ms_product_category` (`product_category_key`),
  ADD CONSTRAINT `FK_ms_product_ms_product_type_product_type_key` FOREIGN KEY (`product_type_key`) REFERENCES `ms_product_type` (`product_type_key`),
  ADD CONSTRAINT `FK_ms_product_ms_risk_profile_risk_profile_key` FOREIGN KEY (`risk_profile_key`) REFERENCES `ms_risk_profile` (`risk_profile_key`);

--
-- Constraints for table `ms_product_bank_account`
--
ALTER TABLE `ms_product_bank_account`
  ADD CONSTRAINT `FK_ms_product_bank_account_ms_bank_account_bank_account_key` FOREIGN KEY (`bank_account_key`) REFERENCES `ms_bank_account` (`bank_account_key`),
  ADD CONSTRAINT `FK_ms_product_bank_account_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`);

--
-- Constraints for table `ms_securities`
--
ALTER TABLE `ms_securities`
  ADD CONSTRAINT `FK_ms_securities_ms_currency_currency_key` FOREIGN KEY (`currency_key`) REFERENCES `ms_currency` (`currency_key`);

--
-- Constraints for table `tr_account`
--
ALTER TABLE `tr_account`
  ADD CONSTRAINT `FK_tr_account_ms_customer_customer_key` FOREIGN KEY (`customer_key`) REFERENCES `ms_customer` (`customer_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_tr_account_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_account_agent`
--
ALTER TABLE `tr_account_agent`
  ADD CONSTRAINT `FK_tr_account_agent_ms_agent_agent_key` FOREIGN KEY (`agent_key`) REFERENCES `ms_agent` (`agent_key`),
  ADD CONSTRAINT `FK_tr_account_agent_ms_branch_branch_key` FOREIGN KEY (`branch_key`) REFERENCES `ms_branch` (`branch_key`),
  ADD CONSTRAINT `FK_tr_account_agent_tr_account_acc_key` FOREIGN KEY (`acc_key`) REFERENCES `tr_account` (`acc_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_balance`
--
ALTER TABLE `tr_balance`
  ADD CONSTRAINT `FK_tr_balance_tr_account_agent_ag_key` FOREIGN KEY (`ag_key`) REFERENCES `tr_account_agent` (`aca_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_currency_rate`
--
ALTER TABLE `tr_currency_rate`
  ADD CONSTRAINT `FK_tr_currency_rate_ms_currency_currency_key` FOREIGN KEY (`currency_key`) REFERENCES `ms_currency` (`currency_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_nav`
--
ALTER TABLE `tr_nav`
  ADD CONSTRAINT `FK_tr_nav_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_prospective_customer`
--
ALTER TABLE `tr_prospective_customer`
  ADD CONSTRAINT `FK_tr_prospective_customer_ms_agent_agent_key` FOREIGN KEY (`agent_key`) REFERENCES `ms_agent` (`agent_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_transaction`
--
ALTER TABLE `tr_transaction`
  ADD CONSTRAINT `FK_tr_transaction_ms_agent_agent_key` FOREIGN KEY (`agent_key`) REFERENCES `ms_agent` (`agent_key`),
  ADD CONSTRAINT `FK_tr_transaction_ms_bank_trans_bank_key` FOREIGN KEY (`trans_bank_key`) REFERENCES `ms_bank` (`bank_key`),
  ADD CONSTRAINT `FK_tr_transaction_ms_branch_branch_key` FOREIGN KEY (`branch_key`) REFERENCES `ms_branch` (`branch_key`),
  ADD CONSTRAINT `FK_tr_transaction_ms_customer_customer_key` FOREIGN KEY (`customer_key`) REFERENCES `ms_customer` (`customer_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_tr_transaction_ms_product_product_key` FOREIGN KEY (`product_key`) REFERENCES `ms_product` (`product_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_tr_transaction_tr_account_agent_aca_key` FOREIGN KEY (`aca_key`) REFERENCES `tr_account_agent` (`aca_key`),
  ADD CONSTRAINT `FK_tr_transaction_tr_transaction_status_trans_status_key` FOREIGN KEY (`trans_status_key`) REFERENCES `tr_transaction_status` (`trans_status_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_tr_transaction_tr_transaction_type_trans_type_key` FOREIGN KEY (`trans_type_key`) REFERENCES `tr_transaction_type` (`trans_type_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_transaction_confirmation`
--
ALTER TABLE `tr_transaction_confirmation`
  ADD CONSTRAINT `FK_tr_transaction_confirmation_tr_transaction_transaction_key` FOREIGN KEY (`transaction_key`) REFERENCES `tr_transaction` (`transaction_key`) ON DELETE CASCADE;

--
-- Constraints for table `tr_transaction_fifo`
--
ALTER TABLE `tr_transaction_fifo`
  ADD CONSTRAINT `FK_tr_transaction_fifo_tr_account_agent_sub_aca_key` FOREIGN KEY (`sub_aca_key`) REFERENCES `tr_account_agent` (`aca_key`),
  ADD CONSTRAINT `FK_tr_transaction_fifo_tr_transaction_trans_red_key` FOREIGN KEY (`trans_red_key`) REFERENCES `tr_transaction` (`transaction_key`),
  ADD CONSTRAINT `FK_tr_transaction_fifo_tr_transaction_trans_sub_key` FOREIGN KEY (`trans_sub_key`) REFERENCES `tr_transaction` (`transaction_key`);

--
-- Constraints for table `tr_transaction_settlement`
--
ALTER TABLE `tr_transaction_settlement`
  ADD CONSTRAINT `FK_tr_transaction_settlement_ms_bank_account_source_bank_accoun~` FOREIGN KEY (`source_bank_account_key`) REFERENCES `ms_bank_account` (`bank_account_key`),
  ADD CONSTRAINT `FK_tr_transaction_settlement_ms_bank_account_target_bank_accoun~` FOREIGN KEY (`target_bank_account_key`) REFERENCES `ms_bank_account` (`bank_account_key`) ON DELETE CASCADE,
  ADD CONSTRAINT `FK_tr_transaction_settlement_tr_transaction_transaction_key` FOREIGN KEY (`transaction_key`) REFERENCES `tr_transaction` (`transaction_key`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

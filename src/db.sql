-- phpMyAdmin SQL Dump
-- version 4.6.5.2
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Apr 28, 2018 at 11:54 AM
-- Server version: 10.1.21-MariaDB
-- PHP Version: 7.0.15

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `medicine`
--

-- --------------------------------------------------------

--
-- Table structure for table `inventory`
--

CREATE TABLE `inventory` (
  `id` int(15) NOT NULL,
  `MedicinelistId` int(15) NOT NULL,
  `Qty` varchar(15) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `medicinelist`
--

CREATE TABLE `medicinelist` (
  `id` int(15) NOT NULL,
  `TradeNames` varchar(50) NOT NULL,
  `GenericNames` varchar(50) NOT NULL,
  `ChemicalNames` varchar(50) NOT NULL,
  `ActivationStatus` varchar(50) NOT NULL,
  `UsedFor` text NOT NULL,
  `Mrp` varchar(15) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `medicinelist`
--

INSERT INTO `medicinelist` (`id`, `TradeNames`, `GenericNames`, `ChemicalNames`, `ActivationStatus`, `UsedFor`, `Mrp`) VALUES
(1, 'Ace ', 'Paracetamol & Caffeine', 'Paracetamol ,  Caffeine', 'true', 'headech', '75'),
(2, 'Ace Plus', 'Paracetamol & Caffeine', 'Paracetamol ,  Caffeine', '', '', '120'),
(3, 'fitamol', 'Paracetamol & Caffeine', 'Paracetamol ,  Caffeine', '', '', '30'),
(4, 'Napa syrup ', 'Paracetamol & Caffeine', 'Paracetamol ,  Caffeine', '', '', '130'),
(8, 'paracitamol', 'Paracetamol & Caffeine', 'Paracetamol ,  Caffeine', '', '', '180'),
(9, 'Napa Extra', 'Paracetamol', '', '', '', '40');

-- --------------------------------------------------------

--
-- Table structure for table `orders`
--

CREATE TABLE `orders` (
  `Id` int(15) NOT NULL,
  `UserId` int(15) NOT NULL,
  `OrderDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `DeliveryDate` varchar(15) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `orders`
--

INSERT INTO `orders` (`Id`, `UserId`, `OrderDate`, `DeliveryDate`) VALUES
(1, 1, '2018-04-17 14:04:48', ''),
(2, 1, '2018-04-17 14:08:35', '');

-- --------------------------------------------------------

--
-- Table structure for table `order_items`
--

CREATE TABLE `order_items` (
  `Id` int(15) NOT NULL,
  `OrderId` int(15) NOT NULL,
  `Productid` int(15) NOT NULL,
  `Qty` varchar(15) NOT NULL,
  `Price` varchar(15) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `order_items`
--

INSERT INTO `order_items` (`Id`, `OrderId`, `Productid`, `Qty`, `Price`) VALUES
(1, 0, 1, '10', ''),
(2, 0, 2, '15', ''),
(3, 0, 3, '20', ''),
(4, 0, 4, '25', ''),
(5, 0, 8, '30', ''),
(6, 1, 2, '15', '40'),
(8, 1, 4, '25', '80'),
(9, 0, 8, '30', ''),
(10, 0, 1, '10', ''),
(11, 0, 8, '30', ''),
(12, 0, 1, '10', ''),
(13, 0, 2, '15', ''),
(14, 0, 4, '25', '50'),
(15, 0, 1, '10', '60'),
(16, 0, 2, '15', '200'),
(17, 0, 1, '1', '60'),
(18, 2, 3, '15', '100'),
(19, 2, 4, '20', '2000'),
(20, 2, 8, '25', '250'),
(21, 2, 1, '5', '100'),
(22, 2, 2, '10', '20');

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `id` int(15) NOT NULL,
  `FirstName` varchar(50) NOT NULL,
  `LastName` varchar(50) NOT NULL,
  `MiddleName` varchar(50) NOT NULL,
  `Gender` varchar(15) NOT NULL,
  `Dob` varchar(15) NOT NULL,
  `Designation` varchar(15) NOT NULL,
  `Area` varchar(50) NOT NULL,
  `Postcode` varchar(15) NOT NULL,
  `Username` varchar(15) NOT NULL,
  `Password` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`id`, `FirstName`, `LastName`, `MiddleName`, `Gender`, `Dob`, `Designation`, `Area`, `Postcode`, `Username`, `Password`) VALUES
(1, 'javed', 'Quayyum', 'Mahmud', 'Male', '22/09/1986', 'Admin', 'Mirpur', '1216', 'javed', '123456'),
(2, 'javed', 'Quayyum', 'mahmud', 'male', '22/09', 'admin', 'Mirpur ', '1216', 'javmah', 'dsvdsvdsvvs'),
(3, 'javed', 'Quayyum', 'mahmud', 'male', '22/09', 'admin', 'Mirpur ', '1216', 'javmah', 'dsvdsvdsvvs'),
(14, 'khaled', 'quayyum', 'mahmud', '', '54544', 'supper Admin', 'Mirpur ', '1216', 'kmq', '123456');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `inventory`
--
ALTER TABLE `inventory`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `medicinelist`
--
ALTER TABLE `medicinelist`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `orders`
--
ALTER TABLE `orders`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `order_items`
--
ALTER TABLE `order_items`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `inventory`
--
ALTER TABLE `inventory`
  MODIFY `id` int(15) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `medicinelist`
--
ALTER TABLE `medicinelist`
  MODIFY `id` int(15) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;
--
-- AUTO_INCREMENT for table `orders`
--
ALTER TABLE `orders`
  MODIFY `Id` int(15) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT for table `order_items`
--
ALTER TABLE `order_items`
  MODIFY `Id` int(15) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=23;
--
-- AUTO_INCREMENT for table `user`
--
ALTER TABLE `user`
  MODIFY `id` int(15) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
